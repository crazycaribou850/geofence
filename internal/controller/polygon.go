package controller

import (
	"io/ioutil"
	"net/http"

	"github.com/geofence/internal/controller/helpers"
	"github.com/geofence/internal/logic"
	"github.com/pquerna/ffjson/ffjson"
	"gopkg.in/go-playground/validator.v9"

	"log"
)

type PolyController struct {
	*helpers.ResponseWritingController
	Validator *validator.Validate
}

func NewPolyController(validator *validator.Validate, log log.Logger) *PolyController {
	return &PolyController{
		ResponseWritingController: &helpers.ResponseWritingController{
			Logger: log,
		},
		Validator: validator,
	}
}

func (c *PolyController) DetermineMembership() func(w http.ResponseWriter, r *http.Request) {
	type IncomingPolyMessage struct {
		Fence *[]logic.Coordinate `json:"fence" validate:"required"`
		Point *logic.Coordinate   `json:"point" validate:"required"`
	}

	type PolyResponse struct {
		Fence    *[]logic.Coordinate `json:"fence"`
		Point    *logic.Coordinate   `json:"point"`
		Position string              `json:"position"`
	}
	return func(w http.ResponseWriter, r *http.Request) {

		body, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			c.Logger.Println("Unprocessable request body", err)
			c.WriteErrorResponse(w, http.StatusInternalServerError, "Could not read body", err)
			return
		}

		var params IncomingPolyMessage
		err = ffjson.Unmarshal(body, &params)
		if err != nil {
			c.Logger.Println("Failed to unmarshal IncomingPolyMessage", err)
			c.WriteErrorResponse(w, http.StatusInternalServerError, "Could not unmarshal input", err)
			return
		}

		err = c.Validator.Struct(params)
		if err != nil {
			c.Logger.Println("Unprocessable Request Body", err)
			c.WriteErrorResponse(w, http.StatusUnprocessableEntity, "Invalid Request Body", err)
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
			c.Logger.Println("PolyResponse Marshal failed", err)
			c.WriteErrorResponse(w, http.StatusInternalServerError, "Could not marshal response", err)
			return
		}
		c.WriteResponse(w, http.StatusOK, responseBody)
	}
}
