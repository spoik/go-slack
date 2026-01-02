package channels

import (
	"encoding/json"
	"net/http"
)

func writeJsonResponse(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ChannelsList(w http.ResponseWriter, r *http.Request) {
	channels := []Channel {
		{"Main"},
		{"Help"},
	}

	writeJsonResponse(w, channels)
}
