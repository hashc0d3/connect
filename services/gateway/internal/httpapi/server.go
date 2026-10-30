package httpapi

import (
	"net/http"
	"time"

	"github.com/connect/gateway/internal/authn"
	"github.com/connect/gateway/internal/sms"
	"github.com/connect/gateway/internal/store"
)

type Handler struct {
	Store    *store.Store
	SMS      *sms.Service
	Internal string
}

func New(db *store.Store, service *sms.Service, internal string) http.Handler {
	h := &Handler{Store: db, SMS: service, Internal: internal}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", h.health)
	mux.HandleFunc("GET /s/{code}", h.shortLink)
	mux.HandleFunc("GET /v2/auth", h.auth)
	mux.HandleFunc("POST /v2/auth", h.auth)
	mux.HandleFunc("GET /v2/sms/send", h.send)
	mux.HandleFunc("POST /v2/sms/send", h.send)
	mux.HandleFunc("GET /v2/sms/status", h.status)
	mux.HandleFunc("POST /v2/sms/status", h.status)
	mux.HandleFunc("GET /v2/sms/list", h.listMessages)
	mux.HandleFunc("POST /v2/sms/list", h.listMessages)
	mux.HandleFunc("GET /v2/balance", h.balance)
	mux.HandleFunc("POST /v2/balance", h.balance)
	mux.HandleFunc("GET /v2/balance/add", h.topUp)
	mux.HandleFunc("POST /v2/balance/add", h.topUp)
	mux.HandleFunc("GET /v2/tariffs", h.tariffs)
	mux.HandleFunc("POST /v2/tariffs", h.tariffs)
	mux.HandleFunc("GET /v2/cards", h.cards)
	mux.HandleFunc("POST /v2/cards", h.cards)
	mux.HandleFunc("GET /v2/sign/list", h.listSigns)
	mux.HandleFunc("POST /v2/sign/list", h.listSigns)
	mux.HandleFunc("POST /v2/sign/add", h.addSign)
	mux.HandleFunc("GET /v2/group/list", h.listGroups)
	mux.HandleFunc("POST /v2/group/list", h.listGroups)
	mux.HandleFunc("POST /v2/group/add", h.addGroup)
	mux.HandleFunc("GET /v2/group/add", h.addGroup)
	mux.HandleFunc("POST /v2/group/delete", h.deleteGroup)
	mux.HandleFunc("GET /v2/group/delete", h.deleteGroup)
	mux.HandleFunc("POST /v2/group/delete-all", h.deleteAllGroups)
	mux.HandleFunc("GET /v2/group/delete-all", h.deleteAllGroups)
	mux.HandleFunc("GET /v2/contact/list", h.listContacts)
	mux.HandleFunc("POST /v2/contact/list", h.listContacts)
	mux.HandleFunc("POST /v2/contact/add", h.addContact)
	mux.HandleFunc("GET /v2/contact/add", h.addContact)
	mux.HandleFunc("POST /v2/contact/delete", h.deleteContact)
	mux.HandleFunc("GET /v2/contact/delete", h.deleteContact)
	mux.HandleFunc("POST /v2/contact/delete-all", h.deleteAllContacts)
	mux.HandleFunc("GET /v2/contact/delete-all", h.deleteAllContacts)
	mux.HandleFunc("GET /v2/blacklist/list", h.listBlacklist)
	mux.HandleFunc("POST /v2/blacklist/list", h.listBlacklist)
	mux.HandleFunc("POST /v2/blacklist/add", h.addBlacklist)
	mux.HandleFunc("GET /v2/blacklist/add", h.addBlacklist)
	mux.HandleFunc("POST /v2/blacklist/delete", h.deleteBlacklist)
	mux.HandleFunc("GET /v2/blacklist/delete", h.deleteBlacklist)
	mux.HandleFunc("POST /v2/otp/create", h.createOTP)
	mux.HandleFunc("POST /v2/otp/validate", h.validateOTP)
	mux.HandleFunc("POST /internal/messages/{id}/release", h.release)
	return mux
}

func (h *Handler) user(w http.ResponseWriter, r *http.Request) (store.User, bool) {
	user, err := authn.FromRequest(r, h.Store, h.Internal)
	if err != nil || user == nil {
		fail(w, err)
		return store.User{}, false
	}
	return *user, true
}

func (h *Handler) health(w http.ResponseWriter, r *http.Request) {
	if err := h.Store.Ping(r.Context()); err != nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"status": "down"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) auth(w http.ResponseWriter, r *http.Request) {
	if _, ok := h.user(w, r); !ok {
		return
	}
	ok(w, nil, "Successful authorization.")
}

func unix(t time.Time) int64 {
	if t.IsZero() {
		return 0
	}
	return t.Unix()
}
