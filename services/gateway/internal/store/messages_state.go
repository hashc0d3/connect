package store

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func (s *Store) ChargeAndRelease(ctx context.Context, m Message) error {
	return s.Tx(ctx, func(tx pgx.Tx) error {
		id := m.ID
		if err := s.DebitTx(ctx, tx, m.UserID, m.Cost, &id, "отправка после модерации"); err != nil {
			return err
		}
		tag, err := tx.Exec(ctx, `
			UPDATE messages
			SET upstream_state = 'pending', status = 0, extend_status = 'queue', date_send = now()
			WHERE id = $1 AND upstream_state = 'held'
		`, m.ID)
		if err != nil {
			return err
		}
		if tag.RowsAffected() == 0 {
			return fmt.Errorf("сообщение не ждёт модерации")
		}
		return nil
	})
}

func (s *Store) ClaimDue(ctx context.Context, limit int) ([]Message, error) {
	rows, err := s.pool.Query(ctx, `
		WITH picked AS (
			SELECT id FROM messages
			WHERE upstream_state = 'pending' AND date_send <= now()
			ORDER BY date_send
			LIMIT $1
			FOR UPDATE SKIP LOCKED
		)
		UPDATE messages m SET upstream_state = 'sending'
		FROM picked
		WHERE m.id = picked.id
		RETURNING m.id, m.user_id::text, m.from_sign, m.number, m.text, m.status, m.extend_status, m.channel,
			m.cost::text, m.date_create, m.date_send, m.date_answer, m.upstream_state, m.callback_url
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	list, _, err := scanMessages(rows, 0)
	return list, err
}

func (s *Store) MarkAccepted(ctx context.Context, id int64, clientID, upstreamID string) error {
	_, err := s.pool.Exec(ctx, `
		UPDATE messages
		SET upstream_state = 'sent', status = 3, extend_status = 'sent', upstream_error = '',
			client_message_id = $2, upstream_message_id = $3, date_answer = now()
		WHERE id = $1
	`, id, clientID, upstreamID)
	return err
}

func (s *Store) MarkUpstream(ctx context.Context, id int64, state string, status int, extend, upstreamErr string, answered bool) error {
	if answered {
		_, err := s.pool.Exec(ctx, `
			UPDATE messages
			SET upstream_state = $2, status = $3, extend_status = $4, upstream_error = $5, date_answer = now()
			WHERE id = $1
		`, id, state, status, extend, upstreamErr)
		return err
	}
	_, err := s.pool.Exec(ctx, `
		UPDATE messages
		SET upstream_state = $2, status = $3, extend_status = $4, upstream_error = $5
		WHERE id = $1
	`, id, state, status, extend, upstreamErr)
	return err
}

func (s *Store) ReleaseHeld(ctx context.Context, id int64) error {
	tag, err := s.pool.Exec(ctx, `
		UPDATE messages
		SET upstream_state = 'pending', status = 0, extend_status = 'queue', date_send = now()
		WHERE id = $1 AND upstream_state = 'held'
	`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("сообщение не ждёт модерации")
	}
	return nil
}
