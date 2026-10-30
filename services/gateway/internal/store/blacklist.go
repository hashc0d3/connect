package store

import (
	"context"
	"fmt"
)

func (s *Store) Blacklist(ctx context.Context, userID string) (map[string]struct{}, error) {
	rows, err := s.pool.Query(ctx, `SELECT number FROM blacklist_entries WHERE user_id = $1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := map[string]struct{}{}
	for rows.Next() {
		var n string
		if err := rows.Scan(&n); err != nil {
			return nil, err
		}
		out[n] = struct{}{}
	}
	return out, rows.Err()
}

type BlackItem struct {
	ID     int
	Number string
	Lname  string
	Fname  string
	Sname  string
}

func (s *Store) AddBlacklist(ctx context.Context, userID, number, lname, fname, sname string) (int, error) {
	var id int
	err := s.pool.QueryRow(ctx, `
		INSERT INTO blacklist_entries (user_id, number, lname, fname, sname)
		VALUES ($1,$2,$3,$4,$5) RETURNING id
	`, userID, number, lname, fname, sname).Scan(&id)
	if isUnique(err) {
		return 0, fmt.Errorf("номер уже в чёрном списке")
	}
	return id, err
}

func (s *Store) ListBlacklist(ctx context.Context, userID string) ([]BlackItem, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT id, number, lname, fname, sname FROM blacklist_entries WHERE user_id = $1 ORDER BY id DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []BlackItem
	for rows.Next() {
		var b BlackItem
		if err := rows.Scan(&b.ID, &b.Number, &b.Lname, &b.Fname, &b.Sname); err != nil {
			return nil, err
		}
		out = append(out, b)
	}
	return out, rows.Err()
}

func (s *Store) DeleteBlacklist(ctx context.Context, userID string, id int64) error {
	tag, err := s.pool.Exec(ctx, `DELETE FROM blacklist_entries WHERE user_id = $1 AND id = $2`, userID, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("запись не найдена")
	}
	return nil
}
