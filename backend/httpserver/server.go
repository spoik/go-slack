package httpserver

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"go-slack/channels"

	"github.com/jackc/pgx/v5"
)

const PORT = 8000

func StartNew(ctx context.Context, db *pgx.Conn) {
	mux := New(ctx, db)

	log.Printf("Server starting on port %d...", PORT)

	servePort := fmt.Sprintf(":%d", PORT)
	http.ListenAndServe(servePort, mux)
}

func New(ctx context.Context, db *pgx.Conn) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("GET /channels", channels.NewChannelList(ctx, db))
	return mux
}
