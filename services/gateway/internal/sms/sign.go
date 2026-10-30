package sms

import (
	"context"
	"strings"

	"github.com/connect/gateway/internal/store"
)

func (s *Service) resolveSign(ctx context.Context, userID, name string, rules store.Rules) (string, string, error) {
	signName := strings.TrimSpace(name)
	channel := "FREE_SIGN"
	if signName == "" {
		if rules.RequireSign {
			return "", "", Fail(422, "Укажите имя отправителя")
		}
		return "Connect", channel, nil
	}

	sig, err := s.Store.SignByName(ctx, userID, signName)
	if err != nil {
		return "", "", err
	}
	if sig == nil {
		return "", "", Fail(422, "Имя отправителя не найдено")
	}
	if rules.RequireSign && sig.Status != 1 {
		return "", "", Fail(422, "Имя отправителя ещё не одобрено")
	}
	if sig.Status == 2 {
		return "", "", Fail(422, "Имя отправителя отклонено")
	}
	if sig.Status == 1 {
		channel = sig.Channel
	}
	return sig.Name, channel, nil
}
