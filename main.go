package main

import (
	"fmt"
	"net/http"
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
	mux.HandleFunc("GET /channels", channelsList)
	mux.HandleFunc("/", http.NotFound)
	return mux
}

func channelsList(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Channels"))
}
