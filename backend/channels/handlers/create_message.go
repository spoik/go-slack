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
	// Get the channel id.
	idStr := r.PathValue("id")

	// Check that the channel id is an integer.
	channelId, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid channel id", http.StatusUnprocessableEntity)
		return
	}

	if !cm.channelExists(w, r.Context(), channelId) {
		return
	}

	err = cm.createMessage(w, r, channelId)
	if err != nil {
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (cm createMessage) channelExists(w http.ResponseWriter, c context.Context, channelId int64) bool {
	// Check that the channel exists
	exists, err := cm.queries.ChannelExists(c, channelId)
	if err != nil {
		slog.Error("Failed to query for channel: %w", "err", err)
		internalServerError(w)
		return false
	}

	if !exists {
		http.Error(w, "Channel does not exist", http.StatusNotFound)
		return false
	}

	return true
}

func (cm createMessage) createMessage(w http.ResponseWriter, r *http.Request, channelId int64) error {
	var newMessage CreateMessageRequest
	err := json.NewDecoder(r.Body).Decode(&newMessage)

	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return err
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(newMessage)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return err
	}

	createParams := queries.CreateMessageParams{
		ChannelID: channelId,
		Message:   newMessage.Message,
	}

	_, err = cm.queries.CreateMessage(r.Context(), createParams)
	if err != nil {
		slog.Error("Failed to insert message into database", "err", err)
		internalServerError(w)
		return err
	}

	return nil
}
