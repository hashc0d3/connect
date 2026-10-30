package sms

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"log"
	"strings"

	"github.com/connect/gateway/internal/policy"
	"github.com/connect/gateway/internal/store"
	"github.com/connect/gateway/internal/upstream"
	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"
)

func (s *Service) CreateOTP(ctx context.Context, userID, phone string) (string, error) {
	if err := s.limit(userID, 30); err != nil {
		return "", err
	}
	normalized, err := policy.NormalizePhone(phone)
	if err != nil {
		return "", Fail(422, err.Error())
	}
	cost, err := s.otpPrice(ctx, normalized)
	if err != nil {
		return "", err
	}
	if err := s.charge(ctx, userID, cost, "otp"); err != nil {
		return "", err
	}
	clientID, err := upstream.ClientMessageID()
	if err != nil {
		_ = s.Store.Refund(ctx, userID, cost, 0, "возврат за otp")
		return "", err
	}
	_, code, err := s.Upstream.SendFlashCall(ctx, clientID, normalized)
	if err != nil {
		_ = s.Store.Refund(ctx, userID, cost, 0, "возврат за otp")
		return "", upstreamOTPError(err)
	}
	hash := ""
	if code != "" {
		hash = hashCode(code)
	}
	if err := s.Store.InsertOTP(ctx, userID, normalized, hash, "sent", clientID); err != nil {
		return "", err
	}
	if s.Upstream.Mock && s.Reveal && code != "" {
		return "Otp sent. Тестовый код: " + code, nil
	}
	return "Otp sent.", nil
}

func (s *Service) ValidateOTP(ctx context.Context, userID, phone, code string) (string, error) {
	normalized, err := policy.NormalizePhone(phone)
	if err != nil {
		return "", Fail(422, err.Error())
	}
	code = strings.TrimSpace(code)
	if len(code) < 4 || len(code) > 8 {
		return "", Fail(422, "Код должен содержать от 4 до 8 цифр")
	}
	return s.validateMockOTP(ctx, userID, normalized, code)
}

func (s *Service) validateMockOTP(ctx context.Context, userID, phone, code string) (string, error) {
	hash, err := s.Store.LatestOTPHash(ctx, userID, phone)
	if err != nil {
		return "", err
	}
	got := hashCode(code)
	if hash == "" {
		return "", Fail(422, "Код ещё не получен")
	}
	if subtle.ConstantTimeCompare([]byte(hash), []byte(got)) != 1 {
		_ = s.Store.MarkOTP(ctx, userID, phone, "failed")
		return "", Fail(422, "OTP does not exist")
	}
	_ = s.Store.MarkOTP(ctx, userID, phone, "validated")
	return "OTP is valid", nil
}

func (s *Service) otpPrice(ctx context.Context, phone string) (decimal.Decimal, error) {
	rules, err := s.Store.Rules(ctx)
	if err != nil {
		return decimal.Zero, err
	}
	tariffs, err := s.Store.Tariffs(ctx)
	if err != nil {
		return decimal.Zero, err
	}
	_, operator := policy.GuessOperator(phone)
	factor := decimal.NewFromInt(1).Add(rules.MarkupPercent.Div(decimal.NewFromInt(100)))
	return store.Price(tariffs, "AUTH", operator).Mul(factor).Round(4), nil
}

func (s *Service) charge(ctx context.Context, userID string, amount decimal.Decimal, comment string) error {
	err := s.Store.Tx(ctx, func(tx pgx.Tx) error {
		if err := s.Store.DebitTx(ctx, tx, userID, amount, nil, comment); err != nil {
			if err.Error() == "недостаточно средств" {
				return Fail(402, "Недостаточно средств на балансе")
			}
			return err
		}
		return nil
	})
	var api *APIError
	if errors.As(err, &api) {
		return api
	}
	return err
}

func upstreamOTPError(err error) error {
	var validation *upstream.ValidationError
	if errors.As(err, &validation) {
		return &APIError{HTTP: 422, Message: validation.Message, Fields: validation.Fields}
	}
	log.Printf("easysms otp: %v", err)
	return Fail(502, "Не удалось выполнить запрос OTP")
}

func hashCode(code string) string {
	sum := sha256.Sum256([]byte(code))
	return hex.EncodeToString(sum[:])
}
