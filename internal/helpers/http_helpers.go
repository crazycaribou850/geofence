package helpers

import (
	"github.com/pquerna/ffjson/ffjson"
	"net/http"
)

type ErrorPayload struct {
	Error ErrorDetails `json:"error"`
}

type ErrorDetails struct {
	Message *string `json:"message"`
	Type    string  `json:"type"`
}

func safeError(e error) *string {
	if e == nil {
		return nil
	}

	s := e.Error()
	return &s
}

func WriteResponse(w http.ResponseWriter, status int, payload []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if _, err := w.Write(payload); err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, "Could not write response", err)
	}
}

func WriteErrorResponse(w http.ResponseWriter, status int, message string, responseErr error) {

	payload := ErrorPayload{
		Error: ErrorDetails{
			Type:    message,
			Message: safeError(responseErr),
		},
	}

	b, err := ffjson.Marshal(payload)

	if err != nil {
		WriteResponse(w, http.StatusInternalServerError, []byte(err.Error()))
		return
	}

	WriteResponse(w, status, b)
}
