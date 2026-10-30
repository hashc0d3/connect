package httpapi

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/connect/gateway/internal/authn"
	"github.com/connect/gateway/internal/sms"
)

func writeJSON(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(payload)
}

func ok(w http.ResponseWriter, data any, message any) {
	writeJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"data":    data,
		"message": message,
	})
}

func fail(w http.ResponseWriter, err error) {
	var api *sms.APIError
	if errors.As(err, &api) {
		writeJSON(w, api.HTTP, map[string]any{
			"success": false,
			"data":    nil,
			"message": api.Message,
			"errors":  api.Fields,
		})
		return
	}
	if errors.Is(err, authn.ErrUnauthenticated) {
		writeJSON(w, http.StatusUnauthorized, map[string]any{
			"success": false,
			"data":    nil,
			"message": "Unauthenticated.",
		})
		return
	}
	log.Printf("http: %v", err)
	writeJSON(w, http.StatusInternalServerError, map[string]any{
		"success": false,
		"data":    nil,
		"message": "Внутренняя ошибка",
	})
}
