package handlers

import (
	"net/http"
	"ya41-56/internal/gophermart/customerror"
	"ya41-56/internal/gophermart/models"
	"ya41-56/internal/shared/contextutil"
	"ya41-56/internal/shared/httputil"
	"ya41-56/internal/shared/repositories"
	"ya41-56/internal/shared/response"
)

type BalanceHandler struct {
	Orders     repositories.Repository[models.Order]
	Withdrawal repositories.Repository[models.Withdrawal]
}

func NewBalanceHandler(orderRepo repositories.Repository[models.Order], withdrawalRepo repositories.Repository[models.Withdrawal]) *BalanceHandler {
	return &BalanceHandler{
		Orders:     orderRepo,
		Withdrawal: withdrawalRepo,
	}
}

// Withdraw

type withdrawRequest struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}

// TODO: accumulate the results of calculations (use Balance model)
func (h *BalanceHandler) Withdraw(w http.ResponseWriter, r *http.Request) {
	var req withdrawRequest
	if err := httputil.ParseJSON(r, &req); err != nil {
		response.Error(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	if req.Sum == 0 {
		response.Error(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	userIDStr, ok := contextutil.GetUserID(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}

	// TODO: Begin the transaction

	orders, err := h.Orders.FindManyByField(r.Context(), "user_id", parseID(userIDStr))
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	withdrawals, err := h.Withdrawal.FindManyByField(r.Context(), "user_id", parseID(userIDStr))
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	sumOfAccruals := float64(0.0)
	var currentOrder models.Order
	for _, order := range orders {
		if order.Status == models.OrderStatusProcessed {
			sumOfAccruals += float64(order.Accrual)
		}

		if order.Number == req.Order {
			currentOrder = order
		}
	}

	if currentOrder.Number != req.Order {
		response.Error(w, http.StatusUnprocessableEntity, http.StatusText(http.StatusUnprocessableEntity))
		return
	}

	sumOfWithdrawals := float64(0.0)
	for _, withdrawal := range withdrawals {
		sumOfWithdrawals += float64(withdrawal.Value)
	}

	if req.Sum > sumOfAccruals-sumOfWithdrawals {
		response.Error(w, http.StatusPaymentRequired, customerror.ErrNotEnoughFunds.Error())
		return
	}

	withdrawal := &models.Withdrawal{
		UserID:  parseID(userIDStr),
		OrderID: currentOrder.ID,
		Value:   float32(req.Sum),
	}

	err = h.Withdrawal.Create(r.Context(), withdrawal)
	if err != nil {
		response.JSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	// TODO: Commit the transaction

	response.JSON(w, http.StatusOK, nil)
}
