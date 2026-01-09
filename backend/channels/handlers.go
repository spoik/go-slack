package channels

import (
	"context"
	"encoding/json"
	"go-slack/channels/queries"
	"log/slog"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func writeJsonResponse(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(data)

	if err != nil {
		slog.Error("Unable to encode json", "error", err)
		internalServerError(w)
	}
}

func internalServerError(w http.ResponseWriter) {
	http.Error(w, "Something went wrong", http.StatusInternalServerError)
}

type ChannelList struct {
	queries *queries.Queries
}

func NewChannelList(ctx context.Context, db *pgxpool.Pool) *ChannelList {
	q := queries.New(db)
	return &ChannelList{queries: q}
}

func (c ChannelList) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	channels, err := c.queries.ListChannels(r.Context())

	if err != nil {
		slog.Error("Unable to fetch channels from the database", "error", err)
		internalServerError(w)
		return
	}

	writeJsonResponse(w, channels)
}
