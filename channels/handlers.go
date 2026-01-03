package channels

import (
	"encoding/json"
	"net/http"

	"github.com/uptrace/bun"
)

func writeJsonResponse(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type ChannelList struct {
	DB *bun.DB
}

func (c ChannelList) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	channels := []Channel{
		{1, "Main"},
		{2, "Help"},
	}

	writeJsonResponse(w, channels)
}
