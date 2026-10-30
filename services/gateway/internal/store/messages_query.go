package store

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"
)

func (s *Store) ListMessages(ctx context.Context, userID, number, text string, page, size int) ([]Message, int, error) {
	if page < 1 {
		page = 1
	}
	if size <= 0 || size > 200 {
		size = 50
	}
	where := `WHERE user_id = $1`
	args := []any{userID}
	if number != "" {
		args = append(args, "%"+number+"%")
		where += fmt.Sprintf(" AND number ILIKE $%d", len(args))
	}
	if text != "" {
		args = append(args, "%"+text+"%")
		where += fmt.Sprintf(" AND text ILIKE $%d", len(args))
	}
	var total int
	if err := s.pool.QueryRow(ctx, `SELECT count(*) FROM messages `+where, args...).Scan(&total); err != nil {
		return nil, 0, err
	}
	args = append(args, size, (page-1)*size)
	rows, err := s.pool.Query(ctx, `
		SELECT id, user_id::text, from_sign, number, text, status, extend_status, channel,
			cost::text, date_create, date_send, date_answer, upstream_state, callback_url
		FROM messages `+where+fmt.Sprintf(` ORDER BY id DESC LIMIT $%d OFFSET $%d`, len(args)-1, len(args)), args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	return scanMessages(rows, total)
}

func (s *Store) GetMessage(ctx context.Context, userID string, id int64) (*Message, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT id, user_id::text, from_sign, number, text, status, extend_status, channel,
			cost::text, date_create, date_send, date_answer, upstream_state, callback_url
		FROM messages WHERE user_id = $1 AND id = $2
	`, userID, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	list, _, err := scanMessages(rows, 0)
	if err != nil || len(list) == 0 {
		return nil, err
	}
	return &list[0], nil
}

func (s *Store) GetMessageAny(ctx context.Context, id int64) (*Message, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT id, user_id::text, from_sign, number, text, status, extend_status, channel,
			cost::text, date_create, date_send, date_answer, upstream_state, callback_url
		FROM messages WHERE id = $1
	`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	list, _, err := scanMessages(rows, 0)
	if err != nil || len(list) == 0 {
		return nil, err
	}
	return &list[0], nil
}

func scanMessages(rows pgx.Rows, total int) ([]Message, int, error) {
	var out []Message
	for rows.Next() {
		var m Message
		var cost string
		if err := rows.Scan(&m.ID, &m.UserID, &m.From, &m.Number, &m.Text, &m.Status, &m.ExtendStatus, &m.Channel, &cost, &m.DateCreate, &m.DateSend, &m.DateAnswer, &m.Upstream, &m.CallbackURL); err != nil {
			return nil, 0, err
		}
		var err error
		m.Cost, err = decimal.NewFromString(cost)
		if err != nil {
			return nil, 0, err
		}
		out = append(out, m)
	}
	return out, total, rows.Err()
}
