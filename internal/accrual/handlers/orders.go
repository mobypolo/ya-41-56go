package handlers

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
	"net/http"
	"ya41-56/internal/accrual/models"
	"ya41-56/internal/accrual/worker"
	"ya41-56/internal/shared/httputil"
	"ya41-56/internal/shared/luhn"
	"ya41-56/internal/shared/repositories"
	"ya41-56/internal/shared/response"
)

type OrdersHandler struct {
	Repo      repositories.Repository[models.Order]
	Processor *worker.Pool
}

func NewOrdersHandler(repo repositories.Repository[models.Order], processor *worker.Pool) *OrdersHandler {
	return &OrdersHandler{
		Repo:      repo,
		Processor: processor,
	}
}

func (h *OrdersHandler) GetOrderByNumber(w http.ResponseWriter, r *http.Request) {
	orderNumber := chi.URLParam(r, "number")
	if orderNumber == "" || !luhn.IsValidLuhn(orderNumber) {
		response.Error(w, http.StatusBadRequest, "missing order number")
		return
	}

	if orderInstance, err := h.Repo.FindByFieldWithPreloads(context.Background(), "number", orderNumber, "Goods"); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Error(w, http.StatusNotFound, "order not found")
		} else {
			response.Error(w, http.StatusInternalServerError, err.Error())
		}
		return
	} else {
		response.JSON(w, http.StatusOK, map[string]interface{}{
			"order":   orderInstance.Number,
			"status":  orderInstance.Status,
			"accrual": orderInstance.Accrual,
		})
	}
}

type CreateOrderRequest struct {
	Order string        `json:"order"`
	Goods []GoodItemDTO `json:"goods"`
}

type GoodItemDTO struct {
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

func (h *OrdersHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req CreateOrderRequest
	if err := httputil.ParseJSON(r, &req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid json")
		return
	}

	if !luhn.IsValidLuhn(req.Order) {
		response.Error(w, http.StatusBadRequest, "invalid order number")
		return
	}

	if len(req.Goods) == 0 {
		response.Error(w, http.StatusBadRequest, "empty goods")
		return
	}

	ctx := r.Context()

	existing, err := h.Repo.FindByField(ctx, "number", req.Order)
	if err == nil && existing != nil {
		response.Error(w, http.StatusConflict, "order already exists")
		return
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	order := h.ConvertToModel(req)

	if err := h.Repo.Create(ctx, order); err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.Processor.Add(worker.OrderTask{OrderID: order.ID})

	response.JSON(w, http.StatusAccepted, map[string]string{
		"order":  order.Number,
		"status": order.Status,
	})
}

func (h *OrdersHandler) ConvertToModel(dto CreateOrderRequest) *models.Order {
	order := &models.Order{
		Number: dto.Order,
		Goods:  make([]models.Good, len(dto.Goods)),
	}

	for i, g := range dto.Goods {
		order.Goods[i] = models.Good{
			Description: g.Description,
			Price:       g.Price,
		}
	}

	return order
}
