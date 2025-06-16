package response

import (
	"encoding/json"
	"net/http"
)

type StatusCode = int

func JSON(w http.ResponseWriter, status int, payload ...interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if len(payload) == 0 || payload[0] == nil {
		return
	}

	if err := json.NewEncoder(w).Encode(payload[0]); err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
	}
}

func Error(w http.ResponseWriter, status StatusCode, message string) {
	JSON(w, status, map[string]string{"error": message})
}
