package channels

import "net/http"

func ChannelsList(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Channels"))
}
