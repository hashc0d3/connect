package store

import (
	"context"
	"fmt"
)

type Group struct {
	ID   int
	Name string
}

func (s *Store) ListGroups(ctx context.Context, userID string, page, size int) ([]Group, int, error) {
	if page < 1 {
		page = 1
	}
	var total int
	if err := s.pool.QueryRow(ctx, `SELECT count(*) FROM "groups" WHERE user_id = $1`, userID).Scan(&total); err != nil {
		return nil, 0, err
	}
	rows, err := s.pool.Query(ctx, `
		SELECT id, name FROM "groups" WHERE user_id = $1 ORDER BY id LIMIT $2 OFFSET $3
	`, userID, size, (page-1)*size)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var out []Group
	for rows.Next() {
		var g Group
		if err := rows.Scan(&g.ID, &g.Name); err != nil {
			return nil, 0, err
		}
		out = append(out, g)
	}
	return out, total, rows.Err()
}

func (s *Store) AddGroup(ctx context.Context, userID, name string) (int, error) {
	var id int
	err := s.pool.QueryRow(ctx, `
		INSERT INTO "groups" (user_id, name) VALUES ($1,$2) RETURNING id
	`, userID, name).Scan(&id)
	if isUnique(err) {
		return 0, fmt.Errorf("группа уже существует")
	}
	return id, err
}

func (s *Store) DeleteGroup(ctx context.Context, userID string, id int64) error {
	tag, err := s.pool.Exec(ctx, `DELETE FROM "groups" WHERE user_id = $1 AND id = $2`, userID, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("группа не найдена")
	}
	return nil
}

func (s *Store) DeleteAllGroups(ctx context.Context, userID string) error {
	_, err := s.pool.Exec(ctx, `DELETE FROM "groups" WHERE user_id = $1`, userID)
	return err
}
