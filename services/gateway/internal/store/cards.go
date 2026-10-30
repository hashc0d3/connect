package store

import "context"

type Card struct {
	ID    int
	Brand string
	Last4 string
}

func (s *Store) Cards(ctx context.Context, userID string) ([]Card, error) {
	rows, err := s.pool.Query(ctx, `SELECT id, brand, last4 FROM cards WHERE user_id = $1 ORDER BY id`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Card
	for rows.Next() {
		var c Card
		if err := rows.Scan(&c.ID, &c.Brand, &c.Last4); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

func (s *Store) CardOwned(ctx context.Context, userID string, id int64) (bool, error) {
	var ok bool
	err := s.pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM cards WHERE user_id = $1 AND id = $2)`, userID, id).Scan(&ok)
	return ok, err
}
