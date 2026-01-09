package main

import (
	"context"
	"go-slack/config"
	"go-slack/httpserver"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cfg := initConfig()
	ctx := context.Background()

	db := initDb(ctx, cfg)
	defer db.Close()

	server := httpserver.NewServer(ctx, db, cfg.Port)

	go func() {
		slog.Info("Server starting", "port", cfg.Port)
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			slog.Error("Server ListenAndServe failed", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("Server Shutdown Failed", "error", err)
		os.Exit(1)
	}
	slog.Info("Server exited properly")
}

func initConfig() *config.Config {
	cfg, err := config.New()
	if err != nil {
		slog.Error("Failed to load config", "error", err)
		os.Exit(1)
	}
	return cfg
}

func initDb(ctx context.Context, cfg *config.Config) *pgxpool.Pool {
	db, err := pgxpool.New(ctx, cfg.DB_URL)

	if err != nil {
		slog.Error("Failed to connect to the database", "error", err)
		os.Exit(1)
	}

	return db
}
