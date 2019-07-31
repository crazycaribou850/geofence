package controller

import (
	helpers2 "github.com/geofence/internal/helpers"
	"github.com/geofence/internal/logic"
	"github.com/pquerna/ffjson/ffjson"
	"gopkg.in/go-playground/validator.v9"
	"io/ioutil"
	"log"
	"net/http"
)

type CircleController struct {
	*helpers2.ResponseWritingController
	Validator *validator.Validate
}

func NewCircleController(validator *validator.Validate, log log.Logger) *CircleController {
	return &CircleController{
		ResponseWritingController: &helpers2.ResponseWritingController{
			Logger: log,
		},
		Validator: validator,
	}
}

func (c *CircleController) DetermineMembership() func(w http.ResponseWriter, r *http.Request) {
	type IncomingCircleMessage struct {
		Fence *logic.RadialFence `json:"fence" validate:"required"`
		Point *[2]float64  `json:"point" validate:"required"`
	}

	type CircleResponse struct {
		Fence    *logic.RadialFence `json:"fence"`
		Point    *[2]float64  `json:"point"`
		Position string             `json:"position"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			c.Logger.Println("Unprocessable request body", err)
			c.WriteErrorResponse(w, http.StatusInternalServerError, "Could not read body", err)
			return
		}

		var params IncomingCircleMessage
		err = ffjson.Unmarshal(body, &params)
		if err != nil {
			c.Logger.Println("Failed to unmarshal IncomingCircleMessage", err)
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
			c.Logger.Println("CircleResponse Marshal Failed", err)
			c.WriteErrorResponse(w, http.StatusInternalServerError, "Could not marshal response", err)
			return
		}
		c.WriteResponse(w, http.StatusOK, responseBody)
	}
}

