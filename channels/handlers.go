package channels

import (
	"context"
	"encoding/json"
	"go-slack/channels/queries"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5"
)

func writeJsonResponse(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(data)

	if err != nil {
		log.Println("Unable to encode json:", err.Error())
		genericInternalServerError(w)
	}
}

func genericInternalServerError(w http.ResponseWriter) {
	http.Error(w, "Something went wrong", http.StatusInternalServerError)
}

type ChannelList struct {
	queries *queries.Queries
}

func NewChannelList(ctx context.Context, db *pgx.Conn) *ChannelList {
	q := queries.New(db)
	return &ChannelList{queries: q}
}

func (c ChannelList) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	channels, err := c.queries.ListChannels(r.Context())

	if err != nil {
		log.Println("Unable to fetch channels from the database:", err.Error())
		genericInternalServerError(w)
		return
	}

	writeJsonResponse(w, channels)
}
