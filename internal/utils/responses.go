package utils

import (
	"net/http"
)

type Response struct {
	Status     string      `json:"status"`
	Payload    interface{} `json:"payload,omitempty"`
	StatusCode int         `json:"status_code"`
}

const (
	ErrorStatus   = "error"
	SuccessStatus = "ok"
)

func Error(w http.ResponseWriter, text interface{}, code int) *Response {
	return respWithStatus(w, text, code, ErrorStatus)
}

func InternalError(w http.ResponseWriter) *Response {
	return respWithStatus(w, "Internal Server Error", 500, ErrorStatus)
}

func Success(w http.ResponseWriter, text interface{}, code int) *Response {
	return respWithStatus(w, text, code, SuccessStatus)
}

func respWithStatus(w http.ResponseWriter, payload interface{}, code int, status string) *Response {

	w.WriteHeader(code)

	return &Response{
		Status:     status,
		Payload:    payload,
		StatusCode: code,
	}
}
