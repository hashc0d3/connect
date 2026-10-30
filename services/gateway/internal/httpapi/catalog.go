package httpapi

import (
	"net/http"
	"strings"

	"github.com/connect/gateway/internal/sms"
	"github.com/shopspring/decimal"
)

func errInvalid(message string) error {
	return sms.Fail(http.StatusUnprocessableEntity, message)
}

func (h *Handler) tariffs(w http.ResponseWriter, r *http.Request) {
	if _, okUser := h.user(w, r); !okUser {
		return
	}
	rules, err := h.Store.Rules(r.Context())
	if err != nil {
		fail(w, err)
		return
	}
	raw, err := h.Store.Tariffs(r.Context())
	if err != nil {
		fail(w, err)
		return
	}
	factor := decimal.NewFromInt(1).Add(rules.MarkupPercent.Div(decimal.NewFromInt(100)))
	grouped := map[string]map[string]string{}
	for key, price := range raw {
		channel, operator, _ := strings.Cut(key, "|")
		if grouped[channel] == nil {
			grouped[channel] = map[string]string{}
		}
		grouped[channel][operator] = price.Mul(factor).Round(2).StringFixed(2)
	}
	ok(w, grouped, nil)
}

func (h *Handler) listSigns(w http.ResponseWriter, r *http.Request) {
	user, okUser := h.user(w, r)
	if !okUser {
		return
	}
	rows, err := h.Store.ListSigns(r.Context(), user.ID)
	if err != nil {
		fail(w, err)
		return
	}
	out := make([]map[string]any, 0, len(rows))
	for _, row := range rows {
		out = append(out, map[string]any{
			"id": row.ID, "name": row.Name, "status": row.Status,
			"extendStatus": row.ExtendStatus, "rejectReason": row.RejectReason, "channel": row.Channel,
		})
	}
	ok(w, map[string]any{"totalCount": len(rows), "items": out}, nil)
}

type signBody struct {
	Name    string `json:"name"`
	Channel string `json:"channel"`
}

func (h *Handler) addSign(w http.ResponseWriter, r *http.Request) {
	user, okUser := h.user(w, r)
	if !okUser {
		return
	}
	var body signBody
	readJSON(r, &body)
	name := formOr(r, body.Name, "name")
	channel := formOr(r, body.Channel, "channel")
	if channel == "" {
		channel = "PAY_SIGN"
	}
	if name == "" {
		fail(w, errInvalid("Укажите имя отправителя"))
		return
	}
	id, err := h.Store.AddSign(r.Context(), user.ID, name, channel)
	if err != nil {
		fail(w, errInvalid(err.Error()))
		return
	}
	ok(w, map[string]any{"id": id, "name": name, "status": 0, "extendStatus": "moderation"}, nil)
}
