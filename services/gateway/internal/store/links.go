package store

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

func (s *Store) PutLink(ctx context.Context, code, userID, target string) error {
	_, err := s.pool.Exec(ctx, `INSERT INTO short_links (code, user_id, target) VALUES ($1,$2,$3)`, code, userID, target)
	return err
}

func (s *Store) InsertShortLink(ctx context.Context, tx pgx.Tx, code, userID, target string) error {
	_, err := tx.Exec(ctx, `
		INSERT INTO short_links (code, user_id, target) VALUES ($1,$2,$3)
	`, code, userID, target)
	return err
}

func (s *Store) ShortTarget(ctx context.Context, code string) (string, error) {
	var target string
	err := s.pool.QueryRow(ctx, `SELECT target FROM short_links WHERE code = $1`, code).Scan(&target)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", nil
	}
	return target, err
}
