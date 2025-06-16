package handlers

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"io"
	"net/http"
	"strings"
	"ya41-56/internal/gophermart/models"
	"ya41-56/internal/gophermart/worker"
	"ya41-56/internal/shared/contextutil"
	"ya41-56/internal/shared/luhn"
	"ya41-56/internal/shared/repositories"
	"ya41-56/internal/shared/response"
)

type OrdersHandler struct {
	Orders  repositories.Repository[models.Order]
	Fetcher *worker.FetchPool
}

func NewOrdersHandler(repo repositories.Repository[models.Order], fetcher *worker.FetchPool) *OrdersHandler {
	return &OrdersHandler{
		Orders:  repo,
		Fetcher: fetcher,
	}
}

func (h *OrdersHandler) Upload(w http.ResponseWriter, r *http.Request) {
	userIDStr, ok := contextutil.GetUserID(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "failed to read body")
		return
	}
	number := strings.TrimSpace(string(body))
	if !luhn.IsValidLuhn(number) {
		response.Error(w, http.StatusUnprocessableEntity, "invalid order number")
		return
	}

	existed, err := h.Orders.FindByField(r.Context(), "number", number)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			response.Error(w, http.StatusConflict, "order already exists")
		}
	}
	if existed.ID > 0 {
		if existed.UserID == parseID(userIDStr) {
			response.Error(w, http.StatusOK, "order already exists")
			return
		}
		response.Error(w, http.StatusConflict, "order already exists")
		return
	}

	order := &models.Order{
		UserID: parseID(userIDStr),
		Number: number,
		Status: models.OrderStatusNew,
	}

	err = h.Orders.Create(r.Context(), order)
	h.Fetcher.Add(order)
	if err != nil {
		response.JSON(w, http.StatusOK, err.Error())
		return
	}

	response.JSON(w, http.StatusAccepted)
}

func (h *OrdersHandler) List(w http.ResponseWriter, r *http.Request) {
	userIDStr, ok := contextutil.GetUserID(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	orders, err := h.Orders.FindManyByField(r.Context(), "user_id", parseID(userIDStr))
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	if len(orders) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	response.JSON(w, http.StatusOK, orders)
}

func parseID(id string) uint {
	var uid uint
	_, _ = fmt.Sscanf(id, "%d", &uid)
	return uid
}
