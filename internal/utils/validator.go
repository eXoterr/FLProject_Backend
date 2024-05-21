package utils

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

func ValidateRequest(r *http.Request, w http.ResponseWriter, payload interface{}) error {
	err := render.DecodeJSON(r.Body, payload)

	if err != nil {
		render.JSON(w, r, Error(w, "invalid json data", 400))
		return err
	}

	err = validator.New(validator.WithRequiredStructEnabled()).Struct(payload)
	if err != nil {
		errs := err.(validator.ValidationErrors)
		resp := ValidationErrors(errs)
		w.WriteHeader(resp.StatusCode)
		render.JSON(w, r, resp)

		return err
	}

	return nil
}
