package main

import (
	"fmt"
	"net/http"
	"go-slack/channels"
)

const PORT = 8080

func main() {
	mux := createServeMux()

	fmt.Printf("Server starting on port %d...", PORT)

	servePort := fmt.Sprintf(":%d", PORT)
	http.ListenAndServe(servePort, mux)
}

func createServeMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /channels", channels.ChannelsList)
	mux.HandleFunc("/", http.NotFound)
	return mux
}
