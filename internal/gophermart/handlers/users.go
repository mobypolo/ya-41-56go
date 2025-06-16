package handlers

import (
	"net/http"
	"strconv"
	"ya41-56/internal/gophermart/models"
	"ya41-56/internal/gophermart/services"
	"ya41-56/internal/shared/httputil"
	"ya41-56/internal/shared/logger"
	"ya41-56/internal/shared/repositories"
	"ya41-56/internal/shared/response"
)

type UsersHandler struct {
	Auth   *services.AuthService
	Orders repositories.Repository[models.Order]
	Users  repositories.Repository[models.User]
}

func NewUsersHandler(authService *services.AuthService, orderRepo repositories.Repository[models.Order]) *UsersHandler {
	return &UsersHandler{
		Auth:   authService,
		Orders: orderRepo,
		Users:  authService.Users,
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
	Withdrawn int     `json:"withdrawn"`
}

func (h *UsersHandler) Balance(w http.ResponseWriter, r *http.Request) {
	orders, err := h.Orders.FindManyByField(r.Context(), "status", models.OrderStatusProcessed)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	// SELECT SUM(accrual) FROM orders WHERE status = "PROCESSED"
	current := float64(0.0)
	for _, order := range orders {
		current += float64(order.Accrual)
	}

	response.JSON(w, http.StatusOK, balanceResponse{
		Current:   current,
		Withdrawn: 0,
	})
}
