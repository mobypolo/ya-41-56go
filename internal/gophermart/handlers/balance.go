package handlers

import (
	"net/http"
	"ya41-56/internal/shared/response"
)

type BalanceHandler struct {
}

func NewBalanceHandler() *BalanceHandler {
	return &BalanceHandler{}
}

// Withdraw

type withdrawRequest struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}

func (h *BalanceHandler) Withdraw(w http.ResponseWriter, _ *http.Request) {
	response.JSON(w, http.StatusOK, nil)
}
