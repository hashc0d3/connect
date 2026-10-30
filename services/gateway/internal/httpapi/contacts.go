package httpapi

import (
	"net/http"
	"time"

	"github.com/connect/gateway/internal/policy"
	"github.com/connect/gateway/internal/store"
)

type contactBody struct {
	Number   string `json:"number"`
	Lname    string `json:"lname"`
	Fname    string `json:"fname"`
	Sname    string `json:"sname"`
	Sex      string `json:"sex"`
	Param1   string `json:"param1"`
	Param2   string `json:"param2"`
	Param3   string `json:"param3"`
	Birthday int64  `json:"birthday"`
	GroupID  int64  `json:"groupId"`
}

func (h *Handler) listContacts(w http.ResponseWriter, r *http.Request) {
	user, okUser := h.user(w, r)
	if !okUser {
		return
	}
	page := int(queryInt(r, "page"))
	rows, total, err := h.Store.ListContacts(r.Context(), user.ID, r.FormValue("number"), queryInt(r, "groupId"), page, 50)
	if err != nil {
		fail(w, err)
		return
	}
	items := make([]map[string]any, 0, len(rows))
	for _, row := range rows {
		groups := make([]map[string]any, 0, len(row.Groups))
		for _, group := range row.Groups {
			groups = append(groups, map[string]any{"id": group.ID, "name": group.Name})
		}
		items = append(items, map[string]any{
			"id": row.ID, "number": row.Number, "lname": row.Lname, "fname": row.Fname, "sname": row.Sname,
			"param1": row.Param1, "param2": row.Param2, "param3": row.Param3, "sex": row.Sex,
			"operator": row.OperatorCode, "extendOperator": row.OperatorName,
			"hlrStatus": row.HLRStatus, "groups": groups,
		})
	}
	if page < 1 {
		page = 1
	}
	ok(w, map[string]any{"items": items, "total": total, "links": pageLinks("/v2/contact/list", page, total, 50)}, nil)
}

func (h *Handler) addContact(w http.ResponseWriter, r *http.Request) {
	user, okUser := h.user(w, r)
	if !okUser {
		return
	}
	var body contactBody
	readJSON(r, &body)
	number, err := policy.NormalizePhone(formOr(r, body.Number, "number"))
	if err != nil {
		fail(w, errInvalid(err.Error()))
		return
	}
	code, name := policy.GuessOperator(number)
	var birthday *time.Time
	if body.Birthday > 0 {
		value := time.Unix(body.Birthday, 0).UTC()
		birthday = &value
	}
	id, err := h.Store.AddContact(r.Context(), user.ID, store.ContactInput{
		Number: number, Lname: formOr(r, body.Lname, "lname"), Fname: formOr(r, body.Fname, "fname"),
		Sname: formOr(r, body.Sname, "sname"), Sex: formOr(r, body.Sex, "sex"),
		Param1: formOr(r, body.Param1, "param1"), Param2: formOr(r, body.Param2, "param2"), Param3: formOr(r, body.Param3, "param3"),
		Birthday: birthday, OperatorCode: code, OperatorName: name, GroupID: body.GroupID, Personal: body.GroupID == 0,
	})
	if err != nil {
		fail(w, errInvalid(err.Error()))
		return
	}
	ok(w, map[string]any{"id": id, "number": number, "extendOperator": name}, nil)
}

func (h *Handler) deleteContact(w http.ResponseWriter, r *http.Request) {
	user, okUser := h.user(w, r)
	if !okUser {
		return
	}
	if err := h.Store.DeleteContact(r.Context(), user.ID, queryInt(r, "id")); err != nil {
		fail(w, errInvalid(err.Error()))
		return
	}
	ok(w, nil, "Contact delete.")
}

func (h *Handler) deleteAllContacts(w http.ResponseWriter, r *http.Request) {
	user, okUser := h.user(w, r)
	if !okUser {
		return
	}
	if err := h.Store.DeleteAllContacts(r.Context(), user.ID); err != nil {
		fail(w, err)
		return
	}
	ok(w, nil, "Contact delete.")
}
