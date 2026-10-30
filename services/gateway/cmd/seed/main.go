package main

import (
	"context"
	"log"

	"github.com/connect/gateway/internal/config"
	"github.com/connect/gateway/internal/migrate"
	"github.com/connect/gateway/internal/store"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	config.ApplyEnvFile(".env")
	config.ApplyEnvFile("../../.env")
	cfg := config.Load()
	if err := migrate.Up(cfg.DatabaseURL); err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, store.NormalizeDSN(cfg.DatabaseURL))
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()
	if err := seed(ctx, pool); err != nil {
		log.Fatal(err)
	}
	log.Print("demo data ready")
}
