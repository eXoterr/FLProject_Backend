package utils

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

func ValidationErrors(errors validator.ValidationErrors) Response {
	errMessages := []string{}
	for _, e := range errors {
		switch e.ActualTag() {
		case "required":
			errMessages = append(errMessages, fmt.Sprintf("field %s is required", e.Field()))
		case "email":
			errMessages = append(errMessages, fmt.Sprintf("field %s must contain a valid email", e.Field()))
		default:
			errMessages = append(errMessages, fmt.Sprintf("field %s is invalid", e.Field()))
		}
	}

	return Response{
		Status:     ErrorStatus,
		Payload:    strings.Join(errMessages, ", "),
		StatusCode: 400,
	}
}
