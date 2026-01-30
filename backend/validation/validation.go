package validation

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/locales/de"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/non-standard/validators"
	de_translations "github.com/go-playground/validator/v10/translations/de"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

var Validate *validator.Validate
var universalTranslator *ut.UniversalTranslator

func init() {
	en := en.New()
	universalTranslator = ut.New(en, en, de.New())

	Validate = validator.New(validator.WithRequiredStructEnabled())

	err := Validate.RegisterValidation("notblank", validators.NotBlank)
	if err != nil {
		panic(fmt.Errorf("Unable to register notblank validation: %w", err))
	}

	setupEnglishTranslations()
	setupGermanTranslations()
}

func setupGermanTranslations() {
	deTrans, found := universalTranslator.GetTranslator("de")
	if !found {
		panic("Unable to find de translations")
	}

	de_translations.RegisterDefaultTranslations(Validate, deTrans)

	Validate.RegisterTranslation(
		"notblank",
		deTrans,
		func(ut ut.Translator) error {
			return ut.Add("notblank", "Der {0} darf nicht leer sein", true)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("notblank", fe.Field())
			return t
		},
	)
}

func setupEnglishTranslations() {
	enTrans, found := universalTranslator.GetTranslator("en")

	if !found {
		panic("Unable to find en translations")
	}

	en_translations.RegisterDefaultTranslations(Validate, enTrans)

	Validate.RegisterTranslation(
		"notblank",
		enTrans,
		func(ut ut.Translator) error {
			return ut.Add("notblank", "{0} must not be blank", true)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("notblank", fe.Field())
			return t
		},
	)
}

func ValidationErrorsToString(r *http.Request, err error) (string, error) {
	var validateErrs validator.ValidationErrors
	if !errors.As(err, &validateErrs) {
		return "", errors.New("The error must be a validator.ValidationErrors.")
	}

	// TODO: Better parsing of Accept-Language. This doesn't handle if the header has
	// multiple languages. It also doesn't handle language preferences (q=0.9)
	trans, _ := universalTranslator.GetTranslator(r.Header.Get("Accept-Language"))

	errMsgs := make([]string, len(validateErrs))

	for _, error := range validateErrs {
		errMsg := error.Translate(trans) + "."
		errMsgs = append(errMsgs, errMsg)
	}

	errMsg := strings.TrimSpace(strings.Join(errMsgs, " "))
	return errMsg, nil
}
