package httpapi

import "net/http"

type nameBody struct {
	Name string `json:"name"`
}

func (h *Handler) listGroups(w http.ResponseWriter, r *http.Request) {
	user, okUser := h.user(w, r)
	if !okUser {
		return
	}
	page := int(queryInt(r, "page"))
	rows, total, err := h.Store.ListGroups(r.Context(), user.ID, page, 50)
	if err != nil {
		fail(w, err)
		return
	}
	items := make([]map[string]any, 0, len(rows))
	for _, row := range rows {
		items = append(items, map[string]any{"id": row.ID, "name": row.Name})
	}
	if page < 1 {
		page = 1
	}
	ok(w, map[string]any{"items": items, "links": pageLinks("/v2/group/list", page, total, 50)}, nil)
}

func (h *Handler) addGroup(w http.ResponseWriter, r *http.Request) {
	user, okUser := h.user(w, r)
	if !okUser {
		return
	}
	var body nameBody
	readJSON(r, &body)
	name := formOr(r, body.Name, "name")
	if name == "" {
		fail(w, errInvalid("Укажите имя группы"))
		return
	}
	id, err := h.Store.AddGroup(r.Context(), user.ID, name)
	if err != nil {
		fail(w, errInvalid(err.Error()))
		return
	}
	ok(w, map[string]any{"id": id, "name": name}, nil)
}

func (h *Handler) deleteGroup(w http.ResponseWriter, r *http.Request) {
	user, okUser := h.user(w, r)
	if !okUser {
		return
	}
	if err := h.Store.DeleteGroup(r.Context(), user.ID, queryInt(r, "id")); err != nil {
		fail(w, errInvalid(err.Error()))
		return
	}
	ok(w, nil, "Group delete.")
}

func (h *Handler) deleteAllGroups(w http.ResponseWriter, r *http.Request) {
	user, okUser := h.user(w, r)
	if !okUser {
		return
	}
	if err := h.Store.DeleteAllGroups(r.Context(), user.ID); err != nil {
		fail(w, err)
		return
	}
	ok(w, nil, "All groups deleted.")
}
