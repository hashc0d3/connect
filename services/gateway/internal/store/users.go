package store

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"
)

type User struct {
	ID      string
	Email   string
	Role    string
	Balance decimal.Decimal
}

func (s *Store) UserByKeyHash(ctx context.Context, hash string) (*User, int, error) {
	var u User
	var keyID int
	var balance string
	err := s.pool.QueryRow(ctx, `
		SELECT u.id::text, u.email, u.role, u.balance::text, k.id
		FROM api_keys k
		JOIN users u ON u.id = k.user_id
		WHERE k.key_hash = $1
	`, hash).Scan(&u.ID, &u.Email, &u.Role, &balance, &keyID)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, 0, nil
	}
	if err != nil {
		return nil, 0, err
	}
	u.Balance, err = decimal.NewFromString(balance)
	if err != nil {
		return nil, 0, err
	}
	_, _ = s.pool.Exec(ctx, `UPDATE api_keys SET last_used_at = now() WHERE id = $1`, keyID)
	return &u, keyID, nil
}

func (s *Store) UserByID(ctx context.Context, id string) (*User, error) {
	var u User
	var balance string
	err := s.pool.QueryRow(ctx, `
		SELECT id::text, email, role, balance::text FROM users WHERE id = $1
	`, id).Scan(&u.ID, &u.Email, &u.Role, &balance)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	u.Balance, err = decimal.NewFromString(balance)
	return &u, err
}

func (s *Store) HasApprovedProfile(ctx context.Context, userID string) (bool, error) {
	var ok bool
	err := s.pool.QueryRow(ctx, `
		SELECT EXISTS(
			SELECT 1 FROM billing_profiles
			WHERE user_id = $1 AND status = 'approved'
		)
	`, userID).Scan(&ok)
	return ok, err
}
