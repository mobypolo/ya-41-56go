package handlers

import (
	"context"
	"gorm.io/gorm"
	"net/http"
	"time"
	"ya41-56/internal/shared/response"
)

func PingHandler(dbConn *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
		defer cancel()

		dbInstance, err := dbConn.DB()
		if err != nil {
			response.Error(w, http.StatusInternalServerError, err.Error())
			return
		}

		if err := dbInstance.PingContext(ctx); err != nil {
			response.Error(w, http.StatusInternalServerError, err.Error())
			return
		}

		response.JSON(w, http.StatusOK, "pong")
	}
}
