package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"
)

func (s *Store) DebitTx(ctx context.Context, tx pgx.Tx, userID string, amount decimal.Decimal, messageID *int64, comment string) error {
	if amount.LessThanOrEqual(decimal.Zero) {
		return nil
	}
	var balance string
	err := tx.QueryRow(ctx, `
		UPDATE users SET balance = balance - $2::numeric
		WHERE id = $1 AND balance >= $2::numeric
		RETURNING balance::text
	`, userID, amount.StringFixed(4)).Scan(&balance)
	if errors.Is(err, pgx.ErrNoRows) {
		return fmt.Errorf("недостаточно средств")
	}
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `
		INSERT INTO ledger_entries (user_id, amount, kind, message_id, comment)
		VALUES ($1, $2::numeric, 'debit', $3, $4)
	`, userID, amount.Neg().StringFixed(4), messageID, comment)
	return err
}

func (s *Store) Credit(ctx context.Context, userID string, amount decimal.Decimal, kind, comment string) (decimal.Decimal, error) {
	var balance string
	err := s.Tx(ctx, func(tx pgx.Tx) error {
		if err := tx.QueryRow(ctx, `
			UPDATE users SET balance = balance + $2::numeric
			WHERE id = $1
			RETURNING balance::text
		`, userID, amount.StringFixed(4)).Scan(&balance); err != nil {
			return err
		}
		_, err := tx.Exec(ctx, `
			INSERT INTO ledger_entries (user_id, amount, kind, comment)
			VALUES ($1, $2::numeric, $3, $4)
		`, userID, amount.StringFixed(4), kind, comment)
		return err
	})
	if err != nil {
		return decimal.Zero, err
	}
	return decimal.NewFromString(balance)
}

func (s *Store) InsertInvoice(ctx context.Context, userID string, amount, bonus decimal.Decimal, status, comment string) error {
	_, err := s.pool.Exec(ctx, `
		INSERT INTO invoices (user_id, amount, bonus, status, comment)
		VALUES ($1, $2::numeric, $3::numeric, $4, $5)
	`, userID, amount.StringFixed(4), bonus.StringFixed(4), status, comment)
	return err
}

func (s *Store) Refund(ctx context.Context, userID string, amount decimal.Decimal, messageID int64, comment string) error {
	if amount.LessThanOrEqual(decimal.Zero) {
		return nil
	}
	_, err := s.Credit(ctx, userID, amount, "refund", comment)
	return err
}
