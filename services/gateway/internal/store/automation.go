package store

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

type AutoRule struct {
	UserID     string
	TemplateID int
	SignName   string
}

func (s *Store) DueAutoRules(ctx context.Context) ([]AutoRule, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT user_id::text, COALESCE(template_id, 0), sign_name
		FROM auto_rules
		WHERE enabled = true
		  AND template_id IS NOT NULL
		  AND (last_run_on IS NULL OR last_run_on < (now() AT TIME ZONE 'Europe/Moscow')::date)
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []AutoRule
	for rows.Next() {
		var r AutoRule
		if err := rows.Scan(&r.UserID, &r.TemplateID, &r.SignName); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, rows.Err()
}

func (s *Store) TemplateBody(ctx context.Context, userID string, id int) (string, error) {
	var body string
	err := s.pool.QueryRow(ctx, `SELECT body FROM templates WHERE user_id = $1 AND id = $2`, userID, id).Scan(&body)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", nil
	}
	return body, err
}

func (s *Store) BirthdayNumbers(ctx context.Context, userID string) ([]string, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT number FROM contacts
		WHERE user_id = $1 AND birthday IS NOT NULL
		  AND EXTRACT(MONTH FROM birthday) = EXTRACT(MONTH FROM (now() AT TIME ZONE 'Europe/Moscow'))
		  AND EXTRACT(DAY FROM birthday) = EXTRACT(DAY FROM (now() AT TIME ZONE 'Europe/Moscow'))
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []string
	for rows.Next() {
		var n string
		if err := rows.Scan(&n); err != nil {
			return nil, err
		}
		out = append(out, n)
	}
	return out, rows.Err()
}

func (s *Store) MarkAutoRan(ctx context.Context, userID string) error {
	_, err := s.pool.Exec(ctx, `
		UPDATE auto_rules SET last_run_on = (now() AT TIME ZONE 'Europe/Moscow')::date WHERE user_id = $1
	`, userID)
	return err
}
