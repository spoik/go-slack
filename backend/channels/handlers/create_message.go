package handlers

import (
	"context"
	"encoding/json"
	"go-slack/channels/queries"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CreateMessageRequest struct {
	Message string `validate:"required"`
}

type createMessage struct {
	queries *queries.Queries
}

func NewCreateMessage(db *pgxpool.Pool) *createMessage {
	q := queries.New(db)
	return &createMessage{queries: q}
}

func (cm createMessage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cId, err := cm.getChannelId(w, r)
	if err != nil {
		return
	}

	if !cm.channelExists(w, r.Context(), cId) {
		return
	}

	msg, err := cm.createMessage(w, r, cId)
	if err != nil {
		return
	}

	w.WriteHeader(http.StatusCreated)
	err = writeJsonResponse(w, msg)
	if err != nil {
		return
	}
}

func (cm createMessage) getChannelId(w http.ResponseWriter, r *http.Request) (int64, error) {
	idStr := r.PathValue("id")

	cId, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid channel id", http.StatusUnprocessableEntity)
		return 0, err
	}

	return cId, nil
}

func (cm createMessage) channelExists(w http.ResponseWriter, c context.Context, channelId int64) bool {
	// Check that the channel exists
	exists, err := cm.queries.ChannelExists(c, channelId)
	if err != nil {
		slog.Error("Failed to query for channel: %w", "error", err)
		internalServerError(w)
		return false
	}

	if !exists {
		http.Error(w, "Channel does not exist", http.StatusNotFound)
		return false
	}

	return true
}

func (cm createMessage) createMessage(w http.ResponseWriter, r *http.Request, channelId int64) (*queries.Message, error) {
	var newMsg CreateMessageRequest
	err := json.NewDecoder(r.Body).Decode(&newMsg)

	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return nil, err
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(newMsg)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)

		return nil, err
	}

	createParams := queries.CreateMessageParams{
		ChannelID: channelId,
		Message:   newMsg.Message,
	}

	msg, err := cm.queries.CreateMessage(r.Context(), createParams)
	if err != nil {
		slog.Error("Failed to insert message into database", "error", err)
		internalServerError(w)
		return nil, err
	}

	return &msg, nil
}
