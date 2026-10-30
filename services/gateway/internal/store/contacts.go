package store

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

type ContactInput struct {
	Number       string
	Lname        string
	Fname        string
	Sname        string
	Sex          string
	Param1       string
	Param2       string
	Param3       string
	Birthday     *time.Time
	OperatorCode int
	OperatorName string
	GroupID      int64
	Personal     bool
}

func (s *Store) AddContact(ctx context.Context, userID string, in ContactInput) (int, error) {
	var id int
	err := s.Tx(ctx, func(tx pgx.Tx) error {
		err := tx.QueryRow(ctx, `
			INSERT INTO contacts (
				user_id, number, lname, fname, sname, birthday, sex, param1, param2, param3,
				operator_code, operator_name, personal
			) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)
			RETURNING id
		`, userID, in.Number, in.Lname, in.Fname, in.Sname, in.Birthday, emptySex(in.Sex), in.Param1, in.Param2, in.Param3, in.OperatorCode, in.OperatorName, in.Personal).Scan(&id)
		if err != nil {
			return err
		}
		if in.GroupID > 0 {
			_, err = tx.Exec(ctx, `
				INSERT INTO contact_groups (contact_id, group_id)
				SELECT $1, id FROM "groups" WHERE id = $2 AND user_id = $3
			`, id, in.GroupID, userID)
		}
		return err
	})
	if isUnique(err) {
		return 0, fmt.Errorf("контакт уже существует")
	}
	return id, err
}

func (s *Store) DeleteContact(ctx context.Context, userID string, id int64) error {
	tag, err := s.pool.Exec(ctx, `DELETE FROM contacts WHERE user_id = $1 AND id = $2`, userID, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("контакт не найден")
	}
	return nil
}

func (s *Store) DeleteAllContacts(ctx context.Context, userID string) error {
	_, err := s.pool.Exec(ctx, `DELETE FROM contacts WHERE user_id = $1`, userID)
	return err
}

type ContactRow struct {
	Contact
	Sex       string
	HLRStatus int
	Groups    []Group
}

func (s *Store) ListContacts(ctx context.Context, userID string, number string, groupID int64, page, size int) ([]ContactRow, int, error) {
	if page < 1 {
		page = 1
	}
	where := `WHERE c.user_id = $1`
	args := []any{userID}
	if number != "" {
		args = append(args, "%"+number+"%")
		where += fmt.Sprintf(" AND c.number ILIKE $%d", len(args))
	}
	if groupID > 0 {
		args = append(args, groupID)
		where += fmt.Sprintf(" AND EXISTS (SELECT 1 FROM contact_groups cg WHERE cg.contact_id = c.id AND cg.group_id = $%d)", len(args))
	}
	var total int
	if err := s.pool.QueryRow(ctx, `SELECT count(*) FROM contacts c `+where, args...).Scan(&total); err != nil {
		return nil, 0, err
	}
	args = append(args, size, (page-1)*size)
	rows, err := s.pool.Query(ctx, `
		SELECT c.id, c.number, c.lname, c.fname, c.sname, c.param1, c.param2, c.param3,
			c.operator_code, c.operator_name, c.sex, c.hlr_status
		FROM contacts c `+where+fmt.Sprintf(` ORDER BY c.id DESC LIMIT $%d OFFSET $%d`, len(args)-1, len(args)), args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var out []ContactRow
	var ids []int
	for rows.Next() {
		var c ContactRow
		if err := rows.Scan(&c.ID, &c.Number, &c.Lname, &c.Fname, &c.Sname, &c.Param1, &c.Param2, &c.Param3, &c.OperatorCode, &c.OperatorName, &c.Sex, &c.HLRStatus); err != nil {
			return nil, 0, err
		}
		out = append(out, c)
		ids = append(ids, c.ID)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	if len(ids) == 0 {
		return out, total, nil
	}
	grows, err := s.pool.Query(ctx, `
		SELECT cg.contact_id, g.id, g.name
		FROM contact_groups cg JOIN "groups" g ON g.id = cg.group_id
		WHERE cg.contact_id = ANY($1)
	`, ids)
	if err != nil {
		return nil, 0, err
	}
	defer grows.Close()
	byContact := map[int][]Group{}
	for grows.Next() {
		var cid, gid int
		var name string
		if err := grows.Scan(&cid, &gid, &name); err != nil {
			return nil, 0, err
		}
		byContact[cid] = append(byContact[cid], Group{ID: gid, Name: name})
	}
	for i := range out {
		out[i].Groups = byContact[out[i].ID]
	}
	return out, total, grows.Err()
}

func emptySex(v string) string {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "male", "female":
		return strings.ToLower(v)
	default:
		return "unknown"
	}
}
