package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/connect/gateway/internal/config"
	"github.com/connect/gateway/internal/grpcapi"
	"github.com/connect/gateway/internal/httpapi"
	"github.com/connect/gateway/internal/migrate"
	"github.com/connect/gateway/internal/sms"
	"github.com/connect/gateway/internal/store"
	"github.com/connect/gateway/internal/upstream"
	"github.com/connect/gateway/internal/worker"
)

func main() {
	config.ApplyEnvFile(".env")
	config.ApplyEnvFile("../../.env")
	cfg := config.Load()
	if err := migrate.Up(cfg.DatabaseURL); err != nil {
		log.Fatal(err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	db, err := store.Open(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	service := &sms.Service{
		Store:    db,
		Upstream: upstream.New(cfg.EasySMSURL, cfg.EasySMSConnect, cfg.EasySMSLogin, cfg.EasySMSPassword, cfg.EasySMSWebhook, cfg.EasySMSMock, cfg.RevealOTP),
		Links:    cfg.PublicLink,
		Reveal:   cfg.RevealOTP,
	}
	go worker.Run(ctx, service)

	httpServer := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           httpapi.New(db, service, cfg.InternalToken),
		ReadHeaderTimeout: 5 * time.Second,
	}
	go func() {
		log.Printf("http %s", cfg.HTTPAddr)
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	grpcServer, lis, err := grpcapi.Listen(cfg.GRPCAddr, db, service)
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		log.Printf("grpc %s", cfg.GRPCAddr)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()

	<-ctx.Done()
	shutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = httpServer.Shutdown(shutdown)
	grpcServer.GracefulStop()
}
