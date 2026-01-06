package channels

import (
	"context"
	"encoding/json"
	"go-slack/channels/queries"
	"net/http"

	"github.com/jackc/pgx/v5"
)

func writeJsonResponse(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type ChannelList struct {
	queries *queries.Queries
	ctx context.Context
}

func NewChannelList(ctx context.Context, db *pgx.Conn) *ChannelList {
	q := queries.New(db)
	return &ChannelList{ctx: ctx, queries: q}
}

func (c ChannelList) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	channels, err := c.queries.ListChannels(c.ctx)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJsonResponse(w, channels)
}
