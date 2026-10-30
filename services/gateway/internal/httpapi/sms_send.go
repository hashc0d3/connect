package httpapi

import (
	"net/http"
	"time"

	"github.com/connect/gateway/internal/sms"
)

type sendBody struct {
	Number      string   `json:"number"`
	Numbers     []string `json:"numbers"`
	Receivers   []struct {
		Phone string `json:"phone"`
	} `json:"receivers"`
	Text        string  `json:"text"`
	Sign        string  `json:"sign"`
	DateSend    int64   `json:"dateSend"`
	SmoothUntil int64   `json:"smoothUntil"`
	CallbackURL string  `json:"callbackUrl"`
	ShortLink   bool    `json:"shortLink"`
	Translit    bool    `json:"translit"`
	GroupIDs    []int64 `json:"groupIds"`
}

func (h *Handler) send(w http.ResponseWriter, r *http.Request) {
	user, okUser := h.user(w, r)
	if !okUser {
		return
	}
	var body sendBody
	readJSON(r, &body)
	in := sms.SendInput{
		Text:        formOr(r, body.Text, "text"),
		Sign:        formOr(r, body.Sign, "sign"),
		CallbackURL: formOr(r, body.CallbackURL, "callbackUrl"),
		ShortLink:   body.ShortLink || r.FormValue("shortLink") == "1",
		Translit:    body.Translit || r.FormValue("translit") == "1" || r.FormValue("translit") == "true",
		GroupIDs:    body.GroupIDs,
		Numbers:     append([]string{}, body.Numbers...),
	}
	if number := formOr(r, body.Number, "number"); number != "" {
		in.Numbers = append(in.Numbers, number)
	}
	for _, receiver := range body.Receivers {
		in.Numbers = append(in.Numbers, receiver.Phone)
	}
	in.Numbers = append(in.Numbers, r.Form["numbers"]...)
	in.Numbers = append(in.Numbers, r.Form["numbers[]"]...)
	if body.DateSend == 0 {
		body.DateSend = queryInt(r, "dateSend")
	}
	if body.DateSend > 0 {
		in.DateSend = time.Unix(body.DateSend, 0).UTC()
	}
	if body.SmoothUntil > 0 {
		in.SmoothUntil = time.Unix(body.SmoothUntil, 0).UTC()
	}
	result, err := h.SMS.Send(r.Context(), user.ID, in)
	if err != nil && (result == nil || len(result.Items) == 0) {
		fail(w, err)
		return
	}
	ok(w, map[string]any{
		"messages":            itemsJSON(result.Items),
		"cost":                result.Cost,
		"excludedDuplicates":  result.Excluded,
		"excludedBlacklist":   result.Blacklisted,
		"heldForModeration":   result.Held,
	}, result.Message)
}

func itemsJSON(items []sms.Item) []map[string]any {
	out := make([]map[string]any, 0, len(items))
	for _, item := range items {
		out = append(out, itemJSON(item))
	}
	return out
}

func itemJSON(item sms.Item) map[string]any {
	payload := map[string]any{
		"id":           item.ID,
		"from":         item.From,
		"number":       item.Number,
		"text":         item.Text,
		"status":       item.Status,
		"extendStatus": item.ExtendStatus,
		"channel":      item.Channel,
		"cost":         item.Cost,
		"dateCreate":   unix(item.DateCreate),
		"dateSend":     unix(item.DateSend),
	}
	if item.DateAnswer != nil {
		payload["dateAnswer"] = unix(*item.DateAnswer)
	}
	return payload
}
