package database

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5"
)

func Connect(ctx context.Context) (*pgx.Conn, error) {
	dsn := os.Getenv("DB_URL")

	return pgx.Connect(ctx, dsn)
}
