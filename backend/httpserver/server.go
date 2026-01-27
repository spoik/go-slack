package httpserver

import (
	"context"
	"fmt"
	"go-slack/channels/handlers"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/cors"
)

func NewServer(ctx context.Context, db *pgxpool.Pool, port int) *http.Server {
	mux := NewMux(ctx, db)
	servePort := fmt.Sprintf(":%d", port)

	// cors.Default is inappropriate for production. This should be changed if this code makes
	// it to a production environment.
	handler := cors.Default().Handler(mux)

	return &http.Server{
		Addr:    servePort,
		Handler: handler,
	}
}

func NewMux(ctx context.Context, db *pgxpool.Pool) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("GET /channels", handlers.NewChannelList(db))
	mux.Handle("POST /channels", handlers.NewCreateChannel(db))
	mux.Handle("GET /channels/{id}/messages", handlers.NewMessageList(db))
	mux.Handle("POST /channels/{id}/messages", handlers.NewCreateMessage(db))
	return mux
}
