package httpserver

import (
	"context"
	"fmt"
	"net/http"

	"go-slack/channels"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewServer(ctx context.Context, db *pgxpool.Pool, port int) *http.Server {
	mux := NewMux(ctx, db)
	servePort := fmt.Sprintf(":%d", port)

	return &http.Server{
		Addr:    servePort,
		Handler: mux,
	}
}

func NewMux(ctx context.Context, db *pgxpool.Pool) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("GET /channels", channels.NewChannelList(ctx, db))
	return mux
}
