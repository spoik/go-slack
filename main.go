package main

import (
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.Handle("/channels", channelHandler{})
	mux.Handle("/", http.NotFoundHandler())

	fmt.Println("Server starting on port 8080...")
	http.ListenAndServe(":8080", mux)
}

type channelHandler struct{}

func (h channelHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Channels"))
}
