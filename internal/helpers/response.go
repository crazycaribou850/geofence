package helpers

import (
	"net/http"

	"github.com/geofence/internal/json"
	"log"
)

type ResponseWritingController struct {
	Logger log.Logger
}

type InsertResponse struct {
	Message string
}

func (c *ResponseWritingController) WriteResponse(w http.ResponseWriter, status int, payload []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if _, err := w.Write(payload); err != nil {
		c.Logger.Println("Could not write response", err)
	}
}

func (c *ResponseWritingController) WriteErrorResponse(w http.ResponseWriter, status int, message string, responseErr error) {

	payload := ErrorPayload{
		Error: ErrorDetails{
			Type:    message,
			Message: safeError(responseErr),
		},
	}

	b, err := json.Marshal(payload)

	if err != nil {
		c.Logger.Println("Could not marshal error response payload", err)
		c.WriteResponse(w, http.StatusInternalServerError, []byte(err.Error()))
		return
	}

	c.WriteResponse(w, status, b)
}



