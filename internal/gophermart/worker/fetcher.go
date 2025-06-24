package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/http"
	"time"
	"ya41-56/internal/shared/logger"

	"ya41-56/internal/gophermart/models"
	"ya41-56/internal/shared/repositories"
)

type AccrualFetcher struct {
	orderRepo repositories.Repository[models.Order]
	baseURL   string
	client    *http.Client
	maxTries  int

	jobs chan *models.Order
}

func NewFetcher(repo repositories.Repository[models.Order], baseURL string) *AccrualFetcher {
	return &AccrualFetcher{
		orderRepo: repo,
		baseURL:   baseURL,
		client:    &http.Client{Timeout: 3 * time.Second},
		maxTries:  5,
		jobs:      make(chan *models.Order, 1000),
	}
}

func (f *AccrualFetcher) Add(order *models.Order) {
	f.jobs <- order
}

func (f *AccrualFetcher) RunWorkers(n int) {
	for i := 0; i < n; i++ {
		go func(id int) {
			for order := range f.jobs {
				f.ProcessOrder(context.Background(), order)
			}
			logger.L().Info("worker stopped", zap.Int("id", id))
		}(i)
	}
}

type accrualResponse struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual"`
}

func (f *AccrualFetcher) ProcessOrder(ctx context.Context, order *models.Order) {
	var (
		attempts = 0
		delay    = time.Second
	)

	for attempts < f.maxTries {
		attempts++

		result, err := f.fetchAccrualResponse(ctx, order.Number)
		if err != nil {
			logger.L().Warn("accrual fetch failed", zap.String("order", order.Number), zap.Error(err))
			time.Sleep(delay)
			delay *= 1
			continue
		}

		err = applyAccrualResult(order, result)
		if err != nil {
			logger.L().Error("failed to convert accrual result", zap.Error(err))
			return
		}

		_ = f.orderRepo.Update(ctx, order)
		return
	}

	order.Status = models.OrderStatusFailedFetch
	if err := f.orderRepo.Update(ctx, order); err != nil {
		logger.L().Error("failed to update order", zap.Error(err))
	}
}

func (f *AccrualFetcher) fetchAccrualResponse(ctx context.Context, number string) (*accrualResponse, error) {
	url := fmt.Sprintf("%s/api/orders/%s", f.baseURL, number)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	resp, err := f.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request failed: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			logger.L().Error("failed to close body", zap.Error(err))
		}
	}()

	if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("temporary status code: %d", resp.StatusCode)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading body: %w", err)
	}

	logger.L().Debug("response body", zap.String("body", string(body)))

	var result accrualResponse

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("decoding json: %w", err)
	}

	if result.Status == models.OrderStatusNew {
		return nil, fmt.Errorf("order in working state")
	}

	return &result, nil
}

func applyAccrualResult(order *models.Order, result *accrualResponse) error {
	order.Status = result.Status
	order.Accrual = float32(result.Accrual)
	return nil
}
