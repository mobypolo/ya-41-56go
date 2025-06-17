package handlers

import (
	"math"
	"net/http"
	"strconv"
	"time"
	"ya41-56/internal/gophermart/models"
	"ya41-56/internal/gophermart/services"
	"ya41-56/internal/shared/contextutil"
	"ya41-56/internal/shared/customstrings"
	"ya41-56/internal/shared/httputil"
	"ya41-56/internal/shared/logger"
	"ya41-56/internal/shared/repositories"
	"ya41-56/internal/shared/response"
)

type UsersHandler struct {
	Auth       *services.AuthService
	Orders     repositories.Repository[models.Order]
	Users      repositories.Repository[models.User]
	Withdrawal repositories.Repository[models.Withdrawal]
}

func NewUsersHandler(authService *services.AuthService, orderRepo repositories.Repository[models.Order], withdrawalRepo repositories.Repository[models.Withdrawal]) *UsersHandler {
	return &UsersHandler{
		Auth:       authService,
		Orders:     orderRepo,
		Users:      authService.Users,
		Withdrawal: withdrawalRepo,
	}
}

func (h *UsersHandler) List(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page <= 0 {
		page = 1
	}

	perPage, _ := strconv.Atoi(r.URL.Query().Get("perPage"))
	if perPage <= 0 {
		perPage = 10
	}

	users, hasNext, err := h.Users.FindAllPaginated(r.Context(), repositories.Pagination{
		Page:    page,
		PerPage: perPage,
	})
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"list":    users,
		"hasNext": hasNext,
	})
}

func (h *UsersHandler) GetByID(_ http.ResponseWriter, _ *http.Request) {

}

func (h *UsersHandler) CreateNew(w http.ResponseWriter, r *http.Request) {
	type registerRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		ClientID string `json:"ClientID"`
		Level    int    `json:"UserLevel"`
		Status   int    `json:"status"`
	}

	var payload registerRequest

	if err := httputil.ParseJSON(r, &payload); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid JSON payload")
		logger.L().Info(err.Error())
		return
	}

	createdUser, err := h.Auth.Register(r.Context(), &models.User{
		Login:    payload.Email,
		Password: payload.Password,
		Status:   models.UserStatus(payload.Status),
	})

	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, createdUser)
}

// Balance

type balanceResponse struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

// TODO: accumulate the results of calculations (use Balance model)
func (h *UsersHandler) Balance(w http.ResponseWriter, r *http.Request) {
	userIDStr, ok := contextutil.GetUserID(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}

	orders, err := h.Orders.FindManyByField(r.Context(), "user_id", customstrings.ParseID(userIDStr))
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	withdrawals, err := h.Withdrawal.FindManyByField(r.Context(), "user_id", customstrings.ParseID(userIDStr))
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	sumOfAccruals := float64(0.0)
	for _, order := range orders {
		if order.Status == models.OrderStatusProcessed {
			sumOfAccruals += float64(order.Accrual)
		}
	}

	sumOfWithdrawals := float64(0.0)
	for _, withdrawal := range withdrawals {
		sumOfWithdrawals += float64(withdrawal.Value)
	}

	response.JSON(w, http.StatusOK, balanceResponse{
		// NOTE: Current = SELECT SUM(accrual) FROM orders WHERE status = "PROCESSED AND user_id = ..."
		Current: float64(math.Round((sumOfAccruals-sumOfWithdrawals)*100) / 100),
		// NOTE: Withdrawn = SELECT SUM(value) FROM withdrawals user_id = ...
		Withdrawn: sumOfWithdrawals,
	})
}

// Withdrawals

type withdrawalResponse struct {
	Order       string    `json:"order"`
	Sum         float64   `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"` // RFC3339
}

func (h *UsersHandler) Withdrawals(w http.ResponseWriter, r *http.Request) {
	userIDStr, ok := contextutil.GetUserID(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}

	// SELECT * FROM withdrawals WHERE user_id = ... ORDER BY created_at ASC
	withdrawals, err := h.Withdrawal.FindManyByField(r.Context(), "user_id", customstrings.ParseID(userIDStr))
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	if len(withdrawals) == 0 {
		response.JSON(w, http.StatusNoContent)

		return
	}

	records := make([]withdrawalResponse, 0, len(withdrawals))
	for i := 0; i < len(withdrawals); i++ {
		records = append(records, withdrawalResponse{
			Order:       withdrawals[i].Order,
			Sum:         float64(withdrawals[i].Value),
			ProcessedAt: withdrawals[i].CreatedAt,
		})
	}

	response.JSON(w, http.StatusOK, records)
}
