package controller

import (
	"github.com/geofence/helpers"
	"github.com/geofence/logic"
	"github.com/pquerna/ffjson/ffjson"
	"io/ioutil"
	"net/http"
)

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
		helpers.WriteErrorResponse(w, http.StatusInternalServerError, "Could not read body", err)
		return
	}

	var params IncomingCircleMessage
	err = ffjson.Unmarshal(body, &params)
	if err != nil {
		helpers.WriteErrorResponse(w, http.StatusInternalServerError, "Could not unmarshal input", err)
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
		helpers.WriteErrorResponse(w, http.StatusInternalServerError, "Could not marshal response", err)
		return
	}
	helpers.WriteResponse(w, http.StatusOK, responseBody)

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
		helpers.WriteErrorResponse(w, http.StatusInternalServerError, "Could not read body", err)
		return
	}

	var params IncomingPolyMessage
	err = ffjson.Unmarshal(body, &params)
	if err != nil {
		helpers.WriteErrorResponse(w, http.StatusInternalServerError, "Could not unmarshal input", err)
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
		helpers.WriteErrorResponse(w, http.StatusInternalServerError, "Could not marshal response", err)
		return
	}
	helpers.WriteResponse(w, http.StatusOK, responseBody)
}
