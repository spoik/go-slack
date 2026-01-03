package main

import (
	"fmt"
	"go-slack/channels"
	"log"
	"net/http"

	"go-slack/database"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/extra/bundebug"
)

const PORT = 8080

func main() {
	db := database.Connect()
	defer db.Close()
	// Add query logging for development
	db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	startMux(db)
}

func startMux(db *bun.DB) {
	mux := createServeMux(db)

	log.Printf("Server starting on port %d...", PORT)

	servePort := fmt.Sprintf(":%d", PORT)
	http.ListenAndServe(servePort, mux)
}

func createServeMux(db *bun.DB) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("GET /channels", channels.ChannelList{DB: db})
	mux.HandleFunc("/", http.NotFound)
	return mux
}
