package migrate

import (
	"database/sql"
	"embed"
	"fmt"

	"github.com/connect/gateway/internal/store"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var files embed.FS

func Up(dsn string) error {
	goose.SetBaseFS(files)
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	db, err := sql.Open("pgx", store.NormalizeDSN(dsn))
	if err != nil {
		return err
	}
	defer db.Close()
	if err := goose.Up(db, "migrations"); err != nil {
		return fmt.Errorf("goose: %w", err)
	}
	return nil
}
