package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/go-playground/form/v4"
	"github.com/go-playground/mold/v4/modifiers"
	"github.com/go-playground/validator/v10"
	apperr "github.com/usesnipet/snipet/internal/app-err"
)

var (
	validate    = validator.New()
	conform     = modifiers.New()
	formDecoder = form.NewDecoder()
)

func ParseBody[T any](r *http.Request, v *T) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		if errors.Is(err, io.EOF) {
			return apperr.BadRequest("request body is required")
		}
		return apperr.BadRequest("invalid JSON body: " + err.Error())
	}
	if err := validate.Struct(v); err != nil {
		return apperr.BadRequest(err.Error())
	}

	if err := conform.Struct(r.Context(), v); err != nil {
		return apperr.BadRequest(err.Error())
	}

	return nil
}

func ParseQuery[T any](r *http.Request, v *T) error {
	if err := formDecoder.Decode(v, r.URL.Query()); err != nil {
		return apperr.BadRequest(err.Error())
	}

	if err := validate.Struct(v); err != nil {
		return apperr.BadRequest(err.Error())
	}

	if err := conform.Struct(r.Context(), v); err != nil {
		return apperr.BadRequest(err.Error())
	}

	return nil
}
