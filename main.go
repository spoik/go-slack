package main

import (
	"fmt"
	"go-slack/channels"
	"log"
	"net/http"
)

const PORT = 8080

func main() {
	// db := database.Connect()
	// defer db.Close()
	startMux()
}

func startMux() {
	mux := createServeMux()

	log.Printf("Server starting on port %d...", PORT)

	servePort := fmt.Sprintf(":%d", PORT)
	http.ListenAndServe(servePort, mux)
}

func createServeMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("GET /channels", channels.ChannelList{})
	mux.HandleFunc("/", http.NotFound)
	return mux
}
