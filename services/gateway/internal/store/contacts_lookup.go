package store

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type Contact struct {
	ID           int
	Number       string
	Lname        string
	Fname        string
	Sname        string
	Param1       string
	Param2       string
	Param3       string
	OperatorCode int
	OperatorName string
}

func (s *Store) ContactsByNumbers(ctx context.Context, userID string, numbers []string) (map[string]Contact, error) {
	if len(numbers) == 0 {
		return map[string]Contact{}, nil
	}
	rows, err := s.pool.Query(ctx, `
		SELECT id, number, lname, fname, sname, param1, param2, param3, operator_code, operator_name
		FROM contacts WHERE user_id = $1 AND number = ANY($2)
	`, userID, numbers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanContacts(rows)
}

func (s *Store) ContactsByGroups(ctx context.Context, userID string, groupIDs []int64) ([]Contact, error) {
	if len(groupIDs) == 0 {
		return nil, nil
	}
	rows, err := s.pool.Query(ctx, `
		SELECT DISTINCT c.id, c.number, c.lname, c.fname, c.sname, c.param1, c.param2, c.param3,
			c.operator_code, c.operator_name
		FROM contacts c
		JOIN contact_groups cg ON cg.contact_id = c.id
		JOIN "groups" g ON g.id = cg.group_id
		WHERE c.user_id = $1 AND g.user_id = $1 AND g.id = ANY($2)
	`, userID, groupIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	set, err := scanContacts(rows)
	if err != nil {
		return nil, err
	}
	out := make([]Contact, 0, len(set))
	for _, c := range set {
		out = append(out, c)
	}
	return out, nil
}

func scanContacts(rows pgx.Rows) (map[string]Contact, error) {
	out := map[string]Contact{}
	for rows.Next() {
		var c Contact
		if err := rows.Scan(&c.ID, &c.Number, &c.Lname, &c.Fname, &c.Sname, &c.Param1, &c.Param2, &c.Param3, &c.OperatorCode, &c.OperatorName); err != nil {
			return nil, err
		}
		out[c.Number] = c
	}
	return out, rows.Err()
}
