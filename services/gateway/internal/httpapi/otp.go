package httpapi

import (
	"crypto/subtle"
	"net/http"
	"strconv"

	"github.com/connect/gateway/internal/authn"
)

type otpBody struct {
	Phone string `json:"phone"`
	Code  string `json:"code"`
}

func (h *Handler) createOTP(w http.ResponseWriter, r *http.Request) {
	user, okUser := h.user(w, r)
	if !okUser {
		return
	}
	var body otpBody
	readJSON(r, &body)
	message, err := h.SMS.CreateOTP(r.Context(), user.ID, formOr(r, body.Phone, "phone"))
	if err != nil {
		fail(w, err)
		return
	}
	ok(w, nil, message)
}

func (h *Handler) validateOTP(w http.ResponseWriter, r *http.Request) {
	user, okUser := h.user(w, r)
	if !okUser {
		return
	}
	var body otpBody
	readJSON(r, &body)
	code := body.Code
	if code == "" {
		code = r.FormValue("code")
	}
	message, err := h.SMS.ValidateOTP(r.Context(), user.ID, formOr(r, body.Phone, "phone"), code)
	if err != nil {
		fail(w, err)
		return
	}
	ok(w, nil, message)
}

func (h *Handler) shortLink(w http.ResponseWriter, r *http.Request) {
	target, err := h.Store.ShortTarget(r.Context(), r.PathValue("code"))
	if err != nil || target == "" {
		http.NotFound(w, r)
		return
	}
	http.Redirect(w, r, target, http.StatusFound)
}

func (h *Handler) release(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("X-Internal-Token")
	if subtle.ConstantTimeCompare([]byte(token), []byte(h.Internal)) != 1 {
		fail(w, authn.ErrUnauthenticated)
		return
	}
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err := h.SMS.Release(r.Context(), id); err != nil {
		fail(w, err)
		return
	}
	ok(w, map[string]any{"id": id}, "Сообщение выпущено в отправку")
}
