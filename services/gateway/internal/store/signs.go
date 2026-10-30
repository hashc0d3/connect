package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type Sign struct {
	ID      int
	Name    string
	Channel string
	Status  int
}

func (s *Store) SignByName(ctx context.Context, userID, name string) (*Sign, error) {
	var sig Sign
	err := s.pool.QueryRow(ctx, `
		SELECT id, name, channel, status
		FROM signs WHERE user_id = $1 AND lower(name) = lower($2)
	`, userID, name).Scan(&sig.ID, &sig.Name, &sig.Channel, &sig.Status)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return &sig, err
}

type SignRow struct {
	ID           int
	Name         string
	Status       int
	ExtendStatus string
	RejectReason string
	Channel      string
}

func (s *Store) ListSigns(ctx context.Context, userID string) ([]SignRow, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT id, name, status, extend_status, reject_reason, channel
		FROM signs WHERE user_id = $1 ORDER BY id
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []SignRow
	for rows.Next() {
		var row SignRow
		if err := rows.Scan(&row.ID, &row.Name, &row.Status, &row.ExtendStatus, &row.RejectReason, &row.Channel); err != nil {
			return nil, err
		}
		out = append(out, row)
	}
	return out, rows.Err()
}

func (s *Store) AddSign(ctx context.Context, userID, name, channel string) (int, error) {
	var id int
	err := s.pool.QueryRow(ctx, `
		INSERT INTO signs (user_id, name, channel, status, extend_status)
		VALUES ($1,$2,$3,0,'moderation') RETURNING id
	`, userID, name, channel).Scan(&id)
	if isUnique(err) {
		return 0, fmt.Errorf("подпись уже существует")
	}
	return id, err
}
