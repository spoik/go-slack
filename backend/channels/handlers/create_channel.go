package handlers

import (
	"encoding/json"
	"errors"
	"go-slack/channels/queries"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/non-standard/validators"
	en_translations "github.com/go-playground/validator/v10/translations/en"
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

	err = validateNewChannel(w, newChan)
	if err != nil {
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

func validateNewChannel(w http.ResponseWriter, newChan CreateChannelRequest) error {
	en := en.New()
	uni := ut.New(en, en)
	trans, _ := uni.GetTranslator("en")

	validate := validator.New(validator.WithRequiredStructEnabled())
	en_translations.RegisterDefaultTranslations(validate, trans)
	err := validate.RegisterValidation("notblank", validators.NotBlank)
	if err != nil {
		slog.Error("Unable to register notblank validation", "error", err)
		internalServerError(w)
		return err
	}

	validate.RegisterTranslation("notblank", trans, func(ut ut.Translator) error {
		return ut.Add("notblank", "{0} must not be blank", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("notblank", fe.Field())
		return t
	})

	err = validate.Struct(newChan)
	if err == nil {
		return nil
	}

	var invalidValidationError *validator.InvalidValidationError
	if errors.As(err, &invalidValidationError) {
		slog.Error("")
		internalServerError(w)
		return err
	}

	var validateErrs validator.ValidationErrors
	if errors.As(err, &validateErrs) {
		errMsgs := make([]string, len(validateErrs))
		for _, error := range validateErrs {
			errMsg := error.Translate(trans) + "."
			errMsgs = append(errMsgs, errMsg)
		}

		errMsg := strings.TrimSpace(strings.Join(errMsgs, " "))

		http.Error(w, errMsg, http.StatusUnprocessableEntity)
		return err
	}

	slog.Error("validation failed", "error", err)
	http.Error(w, "Invalid JSON", http.StatusUnprocessableEntity)

	return err
}
