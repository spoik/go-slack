package handlers

import (
	"go-slack/channels/queries"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
)

type MessageList struct {
	queries *queries.Queries
}

func NewMessageList(db *pgxpool.Pool) *MessageList {
	q := queries.New(db)
	return &MessageList{queries: q}

}

func (ml MessageList) ServeHTTP(w http.ResponseWriter, r *http.Request)	{
	id := r.PathValue("id")

	idStr, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, "Invalid channel id", http.StatusUnprocessableEntity)	
		return
	}

	exists, err := ml.queries.ChannelExists(r.Context(), idStr)

	if err != nil {
		internalServerError(w)
		return
	}

	if !exists {
		http.Error(w, "Channel does not exist", http.StatusNotFound)	
		return
	}
}
