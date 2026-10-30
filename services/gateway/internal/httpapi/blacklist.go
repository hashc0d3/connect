package httpapi

import (
	"net/http"

	"github.com/connect/gateway/internal/policy"
)

type blackBody struct {
	Number string `json:"number"`
	Lname  string `json:"lname"`
	Fname  string `json:"fname"`
	Sname  string `json:"sname"`
}

func (h *Handler) listBlacklist(w http.ResponseWriter, r *http.Request) {
	user, okUser := h.user(w, r)
	if !okUser {
		return
	}
	rows, err := h.Store.ListBlacklist(r.Context(), user.ID)
	if err != nil {
		fail(w, err)
		return
	}
	items := make([]map[string]any, 0, len(rows))
	for _, row := range rows {
		items = append(items, map[string]any{
			"id": row.ID, "number": row.Number, "lname": row.Lname, "fname": row.Fname, "sname": row.Sname,
		})
	}
	ok(w, items, nil)
}

func (h *Handler) addBlacklist(w http.ResponseWriter, r *http.Request) {
	user, okUser := h.user(w, r)
	if !okUser {
		return
	}
	var body blackBody
	readJSON(r, &body)
	number, err := policy.NormalizePhone(formOr(r, body.Number, "number"))
	if err != nil {
		fail(w, errInvalid(err.Error()))
		return
	}
	id, err := h.Store.AddBlacklist(r.Context(), user.ID, number, formOr(r, body.Lname, "lname"), formOr(r, body.Fname, "fname"), formOr(r, body.Sname, "sname"))
	if err != nil {
		fail(w, errInvalid(err.Error()))
		return
	}
	ok(w, map[string]any{"id": id, "number": number}, nil)
}

func (h *Handler) deleteBlacklist(w http.ResponseWriter, r *http.Request) {
	user, okUser := h.user(w, r)
	if !okUser {
		return
	}
	if err := h.Store.DeleteBlacklist(r.Context(), user.ID, queryInt(r, "id")); err != nil {
		fail(w, errInvalid(err.Error()))
		return
	}
	ok(w, nil, "Номер удалён из чёрного списка")
}
