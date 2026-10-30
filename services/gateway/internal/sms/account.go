package sms

import (
	"context"
	"strings"

	"github.com/connect/gateway/internal/store"
	"github.com/shopspring/decimal"
)

func (s *Service) Status(ctx context.Context, userID string, id int64) (*Item, error) {
	msg, err := s.Store.GetMessage(ctx, userID, id)
	if err != nil {
		return nil, err
	}
	if msg == nil {
		return nil, Fail(404, "Сообщение не найдено")
	}
	return toItem(*msg), nil
}

func (s *Service) List(ctx context.Context, userID, number, text string, page int) ([]Item, int, error) {
	rows, total, err := s.Store.ListMessages(ctx, userID, number, text, page, 50)
	if err != nil {
		return nil, 0, err
	}
	out := make([]Item, 0, len(rows))
	for _, row := range rows {
		out = append(out, *toItem(row))
	}
	return out, total, nil
}

func (s *Service) Balance(ctx context.Context, userID string) (string, error) {
	user, err := s.Store.UserByID(ctx, userID)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", Fail(404, "Пользователь не найден")
	}
	return user.Balance.StringFixed(2), nil
}

func (s *Service) TopUp(ctx context.Context, userID string, sum decimal.Decimal, promo string, cardID int64) (decimal.Decimal, decimal.Decimal, error) {
	if sum.LessThan(decimal.NewFromInt(1)) || sum.GreaterThan(decimal.NewFromInt(1_000_000)) {
		return decimal.Zero, decimal.Zero, Fail(422, "Сумма пополнения от 1 до 1 000 000")
	}
	if err := s.requireCard(ctx, userID, cardID); err != nil {
		return decimal.Zero, decimal.Zero, err
	}
	rules, err := s.Store.Rules(ctx)
	if err != nil {
		return decimal.Zero, decimal.Zero, err
	}
	bonus := topUpBonus(sum, promo, rules)
	total := sum.Add(bonus)
	if _, err := s.Store.Credit(ctx, userID, total, "credit", "пополнение"); err != nil {
		return decimal.Zero, decimal.Zero, err
	}
	if err := s.Store.InsertInvoice(ctx, userID, sum, bonus, "paid", "Внутреннее пополнение. Эквайринг не подключён."); err != nil {
		return decimal.Zero, decimal.Zero, err
	}
	return sum, bonus, nil
}

func (s *Service) requireCard(ctx context.Context, userID string, cardID int64) error {
	if cardID <= 0 {
		return nil
	}
	ok, err := s.Store.CardOwned(ctx, userID, cardID)
	if err != nil {
		return err
	}
	if !ok {
		return Fail(422, "Карта не найдена")
	}
	return nil
}

func topUpBonus(sum decimal.Decimal, promo string, rules store.Rules) decimal.Decimal {
	codeOK := rules.PromoCode != "" && strings.EqualFold(strings.TrimSpace(promo), rules.PromoCode)
	if !codeOK && sum.LessThan(rules.PromoMin) {
		return decimal.Zero
	}
	return sum.Mul(rules.PromoPercent).Div(decimal.NewFromInt(100)).Round(2)
}

func toItem(m store.Message) *Item {
	return &Item{
		ID: m.ID, From: m.From, Number: m.Number, Text: m.Text, Status: m.Status,
		ExtendStatus: m.ExtendStatus, Channel: m.Channel, Cost: m.Cost.StringFixed(4),
		DateCreate: m.DateCreate, DateSend: m.DateSend, DateAnswer: m.DateAnswer,
	}
}
