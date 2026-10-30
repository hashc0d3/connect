package store

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

func (s *Store) InsertOTP(ctx context.Context, userID, phone, codeHash, status, clientID string) error {
	_, err := s.pool.Exec(ctx, `
		INSERT INTO otp_requests (user_id, phone, code_hash, status, client_message_id) VALUES ($1,$2,$3,$4,$5)
	`, userID, phone, codeHash, status, clientID)
	return err
}

func (s *Store) LatestOTPHash(ctx context.Context, userID, phone string) (string, error) {
	var hash string
	err := s.pool.QueryRow(ctx, `
		SELECT code_hash FROM otp_requests
		WHERE user_id = $1 AND phone = $2 AND status = 'sent' AND created_at > now() - interval '10 minutes'
		ORDER BY id DESC LIMIT 1
	`, userID, phone).Scan(&hash)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", nil
	}
	return hash, err
}

func (s *Store) MarkOTP(ctx context.Context, userID, phone, status string) error {
	_, err := s.pool.Exec(ctx, `
		UPDATE otp_requests SET status = $3, validated_at = CASE WHEN $3 = 'validated' THEN now() ELSE validated_at END
		WHERE id = (
			SELECT id FROM otp_requests
			WHERE user_id = $1 AND phone = $2
			ORDER BY id DESC LIMIT 1
		)
	`, userID, phone, status)
	return err
}
