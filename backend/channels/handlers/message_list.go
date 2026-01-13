package handlers

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"go-slack/channels/queries"
	"net/http"
	"strconv"
)

type MessageList struct {
	queries *queries.Queries
}

func NewMessageList(db *pgxpool.Pool) *MessageList {
	q := queries.New(db)
	return &MessageList{queries: q}

}

func (ml MessageList) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	channelId, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid channel id", http.StatusUnprocessableEntity)
		return
	}

	exists, err := ml.queries.ChannelExists(r.Context(), channelId)

	if err != nil {
		internalServerError(w)
		return
	}

	if !exists {
		http.Error(w, "Channel does not exist", http.StatusNotFound)
		return
	}

	messages, err := ml.queries.MessagesInChannel(r.Context(), channelId)

	if err != nil {
		internalServerError(w)
		return
	}

	if messages == nil {
		messages = []queries.Message{}
	}

	writeJsonResponse(w, messages)
}
