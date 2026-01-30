package handlers

import (
	"encoding/json"
	"errors"
	"go-slack/channels/queries"
	"go-slack/validation"
	"log/slog"
	"net/http"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CreateChannelRequest struct {
	Name string `validate:"required,notblank"`
}

type createChannel struct {
	queries *queries.Queries
}

func NewCreateChannel(db *pgxpool.Pool) *createChannel {
	q := queries.New(db)
	return &createChannel{queries: q}
}

func (cc createChannel) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var newChan CreateChannelRequest
	err := json.NewDecoder(r.Body).Decode(&newChan)

	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if !validateNewChannel(w, r, newChan) {
		return
	}

	createdChan, err := cc.queries.CreateChannel(r.Context(), newChan.Name)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			http.Error(w, "Channel name already taken.", http.StatusUnprocessableEntity)
			return
		}

		slog.Error("Unable to create new channel", "error", err)
		internalServerError(w)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdChan)
}

func validateNewChannel(w http.ResponseWriter, r *http.Request, newChan CreateChannelRequest) bool {
	err := validation.Validate.Struct(newChan)
	if err == nil {
		return true
	}

	errMsg, err := validation.ValidationErrorsToString(r, err)

	if err != nil {
		slog.Error("Failed to convert validation errors to a string", "error", err)
		internalServerError(w)
		return false
	}

	http.Error(w, errMsg, http.StatusUnprocessableEntity)
	return false
}
