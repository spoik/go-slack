package main

import (
	"fmt"
	"net/http"
)

func main() {
	mux := createServeMux()
	fmt.Println("Server starting on port 8080...")
	http.ListenAndServe(":8080", mux)
}

func createServeMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("GET /channels", channelHandler{})
	mux.Handle("/", http.NotFoundHandler())
	return mux
}

type channelHandler struct{}

func (h channelHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Channels"))
}
