package httpapi

import (
	"net/http"

	"github.com/shopspring/decimal"
)

func (h *Handler) balance(w http.ResponseWriter, r *http.Request) {
	user, okUser := h.user(w, r)
	if !okUser {
		return
	}
	value, err := h.SMS.Balance(r.Context(), user.ID)
	if err != nil {
		fail(w, err)
		return
	}
	ok(w, map[string]string{"balance": value}, nil)
}

type topUpBody struct {
	Sum   string `json:"sum"`
	Promo string `json:"promo"`
	Card  int64  `json:"cardId"`
}

func (h *Handler) topUp(w http.ResponseWriter, r *http.Request) {
	user, okUser := h.user(w, r)
	if !okUser {
		return
	}
	var body topUpBody
	readJSON(r, &body)
	raw := formOr(r, body.Sum, "sum")
	sum, err := decimal.NewFromString(raw)
	if err != nil {
		fail(w, errInvalid("Сумма пополнения должна быть числом"))
		return
	}
	if body.Card == 0 {
		body.Card = queryInt(r, "cardId")
	}
	paid, bonus, err := h.SMS.TopUp(r.Context(), user.ID, sum, formOr(r, body.Promo, "promo"), body.Card)
	if err != nil {
		fail(w, err)
		return
	}
	ok(w, map[string]string{"sum": paid.StringFixed(2), "bonus": bonus.StringFixed(2)}, "payment successful")
}

func (h *Handler) cards(w http.ResponseWriter, r *http.Request) {
	user, okUser := h.user(w, r)
	if !okUser {
		return
	}
	cards, err := h.Store.Cards(r.Context(), user.ID)
	if err != nil {
		fail(w, err)
		return
	}
	out := make([]map[string]any, 0, len(cards))
	for _, card := range cards {
		out = append(out, map[string]any{
			"id":     card.ID,
			"number": card.Brand + "*****" + card.Last4,
		})
	}
	ok(w, out, nil)
}
