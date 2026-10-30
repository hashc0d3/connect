package store

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

func isUnique(err error) bool {
	var pg *pgconn.PgError
	return errors.As(err, &pg) && pg.Code == "23505"
}
