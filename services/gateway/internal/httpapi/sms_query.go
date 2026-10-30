package httpapi

import "net/http"

func (h *Handler) status(w http.ResponseWriter, r *http.Request) {
	user, okUser := h.user(w, r)
	if !okUser {
		return
	}
	item, err := h.SMS.Status(r.Context(), user.ID, queryInt(r, "id"))
	if err != nil {
		fail(w, err)
		return
	}
	ok(w, itemJSON(*item), nil)
}

func (h *Handler) listMessages(w http.ResponseWriter, r *http.Request) {
	user, okUser := h.user(w, r)
	if !okUser {
		return
	}
	page := int(queryInt(r, "page"))
	if page < 1 {
		page = 1
	}
	items, total, err := h.SMS.List(r.Context(), user.ID, r.FormValue("number"), r.FormValue("text"), page)
	if err != nil {
		fail(w, err)
		return
	}
	ok(w, map[string]any{
		"messages": itemsJSON(items),
		"total":    total,
		"links":    pageLinks("/v2/sms/list", page, total, 50),
	}, nil)
}
