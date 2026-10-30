package httpapi

import (
	"crypto/subtle"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/connect/gateway/internal/sms"
)

func (h *Handler) easysmsWebhook(w http.ResponseWriter, r *http.Request) {
	raw, err := io.ReadAll(io.LimitReader(r.Body, 1<<20))
	if err != nil {
		fail(w, err)
		return
	}
	var body struct {
		Header struct {
			Token       string `json:"token"`
			WebhookType string `json:"webhook_type"`
		} `json:"Header"`
		Payload struct {
			ClientMessageID string `json:"client_message_id"`
			MessageID       string `json:"message_id"`
			State           string `json:"state"`
			Code            string `json:"code"`
		} `json:"Payload"`
	}
	if json.Unmarshal(raw, &body) != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"message": "Некорректный JSON"})
		return
	}
	token := h.SMS.Upstream.WebhookToken
	if token != "" && subtle.ConstantTimeCompare([]byte(token), []byte(body.Header.Token)) != 1 {
		writeJSON(w, http.StatusForbidden, map[string]string{"message": "Неверный token"})
		return
	}
	if body.Header.WebhookType != "" && body.Header.WebhookType != "status" {
		writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
		return
	}
	if err := h.SMS.ApplyWebhook(r.Context(), sms.WebhookStatus{
		ClientMessageID: strings.TrimSpace(body.Payload.ClientMessageID),
		UpstreamID:      strings.TrimSpace(body.Payload.MessageID),
		State:           strings.TrimSpace(body.Payload.State),
		Code:            strings.TrimSpace(body.Payload.Code),
	}); err != nil {
		fail(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}
