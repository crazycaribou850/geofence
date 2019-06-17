package controller

import (
	"encoding/json"
	"github.com/geofence/logic"
	"github.com/pquerna/ffjson/ffjson"
	"io/ioutil"
	"net/http"
)

type ErrorPayload struct {
	Error ErrorDetails `json:"error"`
}

type ErrorDetails struct {
	Message *string `json:"message"`
	Type    string  `json:"type"`
}

func WriteErrorResponse(w http.ResponseWriter, status int, message string, responseErr error) {

	payload := ErrorPayload{
		Error: ErrorDetails{
			Type:    message,
			Message: safeError(responseErr),
		},
	}

	b, err := json.Marshal(payload)

	if err != nil {
		WriteResponse(w, http.StatusInternalServerError, []byte(err.Error()))
		return
	}

	WriteResponse(w, status, b)
}

func safeError(e error) *string {
	if e == nil {
		return nil
	}

	s := e.Error()
	return &s
}

func CircleHandler(w http.ResponseWriter, r *http.Request) {

	type IncomingCircleMessage struct {
		Fence *logic.RadialFence `json:"fence"`
		Point *logic.Coordinate  `json:"point"`
	}

	type CircleResponse struct {
		Fence    *logic.RadialFence `json:"fence"`
		Point    *logic.Coordinate  `json:"point"`
		Position string             `json:"position"`
	}
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, "Could not read body", err)
		return
	}

	var params IncomingCircleMessage
	err = ffjson.Unmarshal(body, &params)
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, "Could not unmarshal input", err)
		return
	}

	point := params.Point
	fence := params.Fence

	result := logic.InRadius(*point, *fence)
	var position string
	if result {
		position = "Inside"
	} else {
		position = "Outside"
	}
	responseBodyInfo := CircleResponse{fence, point, position}
	responseBody, err := ffjson.Marshal(responseBodyInfo)
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, "Could not marshal response", err)
		return
	}
	WriteResponse(w, http.StatusOK, responseBody)

}

func PolyHandler(w http.ResponseWriter, r *http.Request) {

	type IncomingPolyMessage struct {
		Fence *[]logic.Coordinate `json:"fence"`
		Point   *logic.Coordinate	`json:"point"`
	}

	type PolyResponse struct {
		Fence    *[]logic.Coordinate `json:"fence"`
		Point    *logic.Coordinate  `json:"point"`
		Position string             `json:"position"`
	}

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, "Could not read body", err)
		return
	}

	var params IncomingPolyMessage
	err = ffjson.Unmarshal(body, &params)
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, "Could not unmarshal input", err)
		return
	}

	point := params.Point
	fence := params.Fence

	result := logic.InPoly(*point, *fence)
	var position string
	if result {
		position = "Inside"
	} else {
		position = "Outside"
	}
	responseBodyInfo := PolyResponse{fence, point, position}
	responseBody, err := ffjson.Marshal(responseBodyInfo)
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, "Could not marshal response", err)
		return
	}
	WriteResponse(w, http.StatusOK, responseBody)
}

func WriteResponse(w http.ResponseWriter, status int, payload []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if _, err := w.Write(payload); err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, "Could not write response", err)
	}
}
