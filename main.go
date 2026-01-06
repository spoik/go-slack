package main

import (
	"context"
	"fmt"
	"go-slack/channels"
	"go-slack/database"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5"
)

const PORT = 8080

func main() {
	ctx := context.Background()
	db, err := database.Connect(ctx)

	if err != nil {
		log.Printf("Failed to connect to the database: %e\n", err)
		return
	}
	
	defer db.Close(ctx)
	startMux(ctx, db)
}

func startMux(ctx context.Context, db *pgx.Conn) {
	mux := createServeMux(ctx, db)

	log.Printf("Server starting on port %d...", PORT)

	servePort := fmt.Sprintf(":%d", PORT)
	http.ListenAndServe(servePort, mux)
}

func createServeMux(ctx context.Context, db *pgx.Conn) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("GET /channels", channels.NewChannelList(ctx, db))
	return mux
}
