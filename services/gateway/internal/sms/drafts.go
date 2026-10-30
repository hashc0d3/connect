package sms

import (
	"context"
	"time"

	"github.com/connect/gateway/internal/policy"
	"github.com/connect/gateway/internal/store"
	"github.com/shopspring/decimal"
)

func (s *Service) buildDrafts(ctx context.Context, userID string, in SendInput, rules store.Rules, tariffs map[string]decimal.Decimal, numbers []string, signName, channel string) ([]draft, decimal.Decimal, error) {
	baseText := in.Text
	contacts, err := s.Store.ContactsByNumbers(ctx, userID, numbers)
	if err != nil {
		return nil, decimal.Zero, err
	}
	linkMap, err := s.shortLinks(ctx, userID, baseText, in.ShortLink)
	if err != nil {
		return nil, decimal.Zero, err
	}

	now := time.Now().UTC()
	start := in.DateSend
	if start.IsZero() {
		start = now
	}
	if !in.SmoothUntil.IsZero() && !in.SmoothUntil.After(start) {
		return nil, decimal.Zero, Fail(422, "Окончание плавной отправки должно быть позже начала")
	}
	times := spread(len(numbers), start, in.SmoothUntil)
	factor := decimal.NewFromInt(1).Add(rules.MarkupPercent.Div(decimal.NewFromInt(100)))

	drafts := make([]draft, 0, len(numbers))
	total := decimal.Zero
	for i, number := range numbers {
		text := renderText(baseText, contacts[number], in, rules, signName, linkMap, s.Links)
		if len([]rune(text)) < rules.MinText {
			return nil, decimal.Zero, Fail(422, "Текст для "+number+" короче минимальной длины")
		}
		parts, _ := policy.Analyze(text)
		opName := contacts[number].OperatorName
		if opName == "" || opName == "OTHER" {
			_, opName = policy.GuessOperator(number)
		}
		cost := store.Price(tariffs, channel, opName).Mul(factor).Mul(decimal.NewFromInt(int64(parts))).Round(4)
		total = total.Add(cost)
		drafts = append(drafts, draft{
			Number: number, Text: text, Parts: parts, Cost: cost, Operator: opName, When: times[i],
		})
	}
	return drafts, total, nil
}

func (s *Service) shortLinks(ctx context.Context, userID, text string, enabled bool) (map[string]string, error) {
	links := map[string]string{}
	if !enabled {
		return links, nil
	}
	for _, raw := range policy.ExtractURLs(text) {
		if _, ok := links[raw]; ok {
			continue
		}
		code, err := randomCode(7)
		if err != nil {
			return nil, err
		}
		if err := s.Store.PutLink(ctx, code, userID, raw); err != nil {
			return nil, err
		}
		links[raw] = code
	}
	return links, nil
}

func renderText(base string, contact store.Contact, in SendInput, rules store.Rules, signName string, links map[string]string, publicBase string) string {
	text := policy.Personalize(base, policy.Contact{
		Number: contact.Number, Lname: contact.Lname, Fname: contact.Fname, Sname: contact.Sname,
		Param1: contact.Param1, Param2: contact.Param2, Param3: contact.Param3, OperatorName: contact.OperatorName,
	})
	if in.Translit {
		text = policy.Transliterate(text)
	}
	if in.ShortLink {
		text = policy.Shorten(text, publicBase, links)
	}
	if rules.PrefixSign && signName != "" {
		text = signName + ": " + text
	}
	return text
}
