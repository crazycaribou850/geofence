package controller

import (
	"github.com/geofence/internal/helpers"
	"github.com/geofence/internal/model"
	"github.com/geofence/internal/repository"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/geofence/internal/json"
	"github.com/geofence/internal/logic"
	"gopkg.in/go-playground/validator.v9"

	"log"
)

type PolyController struct {
	*helpers.ResponseWritingController
	Validator *validator.Validate
	Repository repository.PolygonPostgresRepository
}

type IncomingFindClosestRequest struct {
	StoreID int `json:"store_id"`
	Point *model.PointGeometry `json:"point" validate:"required"`
}

type IncomingFindEnclosingRequest struct {
	StoreID int `json:"store_id" validate:"required"`
	MetroID int `json:"metro_id" validate:"required"`
	ZoneID int `json:"zone_id" validate:"required"`
	Point *model.PointGeometry `json:"point" validate:"required"`
}

func NewPolyController(validator *validator.Validate, log log.Logger, db *sqlx.DB) *PolyController {
	return &PolyController{
		ResponseWritingController: &helpers.ResponseWritingController{
			Logger: log,
		},
		Validator: validator,
		Repository: repository.PolygonPostgresRepository{DB: *db},
	}
}

func (c *PolyController) DetermineMembership() func(w http.ResponseWriter, r *http.Request) {
	type IncomingMessage struct {
		Geom *model.PolyGeometry `json:"geom" validate:"required"`
		Point *[2]float64 `json:"point" validate:"required"`
	}

	type PolyResponse struct {
		Geom    *model.PolyGeometry `json:"geom"`
		Point    *[2]float64   `json:"point"`
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

		var params IncomingMessage
		err = json.Unmarshal(body, &params)
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
		geom := params.Geom

		result := logic.InPoly(*point, geom.Coordinates[0])
		var position string
		if result {
			position = "Inside"
		} else {
			position = "Outside"
		}
		responseBodyInfo := PolyResponse{geom, point, position}
		responseBody, err := json.Marshal(responseBodyInfo)
		if err != nil {
			c.Logger.Println("PolyResponse Marshal failed", err)
			c.WriteErrorResponse(w, http.StatusInternalServerError, "Could not marshal response", err)
			return
		}
		c.WriteResponse(w, http.StatusOK, responseBody)
	}
}

func (c PolyController) InsertPolygon() func(w http.ResponseWriter, r *http.Request) {

	type IncomingPolygon struct {
		ID      int                `json:"id"`
		Polygon model.PolyGeometry `json:"polygon"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			c.Logger.Println("Unprocessable request body", err)
			c.WriteErrorResponse(w, http.StatusInternalServerError, "Could not read body", err)
			return
		}

		var params IncomingPolygon
		err = json.Unmarshal(body, &params)
		if err != nil {
			c.Logger.Println("Failed to unmarshal IncomingPolyLocation", err)
			c.WriteErrorResponse(w, http.StatusInternalServerError, "Could not unmarshal input", err)
			return
		}
		err = c.Validator.Struct(params)
		if err != nil {
			c.Logger.Println("Unprocessable Request Body", err)
			c.WriteErrorResponse(w, http.StatusUnprocessableEntity, "Invalid Request Body", err)
			return
		}
		err = c.Repository.InsertPolygon(params.ID, params.Polygon)
		if err != nil {
			c.Logger.Println("Failed to insert into table")
			c.WriteErrorResponse(w, http.StatusUnprocessableEntity, "Invalid Insert Request", err)
			return
		}
		result := helpers.InsertResponse{"Insert Success!"}
		responseBody, err := json.Marshal(result)
		if err != nil {
			c.Logger.Println("Response Marshal failed", err)
			c.WriteErrorResponse(w, http.StatusInternalServerError, "Could not marshal response", err)
			return
		}
		c.WriteResponse(w, http.StatusOK, responseBody)
	}
}



func (c PolyController) Ping() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		result, err := c.Repository.GetAll()
		if err != nil {
			c.Logger.Println("Failed to get all from table")
			c.WriteErrorResponse(w, http.StatusUnprocessableEntity, "Invalid get all Request", err)
			return
		}
		responseBody, err := json.Marshal(result)
		if err != nil {
			c.Logger.Println("PolyRow Marshal failed", err)
			c.WriteErrorResponse(w, http.StatusInternalServerError, "Could not marshal response", err)
			return
		}
		c.WriteResponse(w, http.StatusOK, responseBody)
	}
}

func (c *PolyController) DetermineGeogMembership() func(w http.ResponseWriter, r *http.Request) {
	type IncomingMessage struct {
		Geom *model.PolyGeometry `json:"geom" validate:"required"`
		Point *model.PointGeometry `json:"point" validate:"required"`
	}

	type PolyResponse struct {
		Geom    *model.PolyGeometry `json:"geom"`
		Point    *model.PointGeometry  `json:"point"`
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

		var params IncomingMessage
		err = json.Unmarshal(body, &params)
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
		geom := params.Geom

		geomJSON, err:= json.Marshal(geom)
		if err != nil {
			c.Logger.Println("Failed to Marshal geomJSON object")
			c.WriteErrorResponse(w, http.StatusInternalServerError, "Could not marshal geomJSON", err)
		}
		geomString := string(geomJSON)

		pointJSON, err := json.Marshal(point)
		if err != nil {
			c.Logger.Println("Failed to Marshal pointJSON object")
			c.WriteErrorResponse(w, http.StatusInternalServerError, "Could not marshal geomJSON", err)
		}
		pointString := string(pointJSON)


		result, err := c.Repository.Intersects(geomString, pointString)
		if err != nil {
			c.Logger.Println("DB Query failed")
			c.WriteErrorResponse(w, http.StatusInternalServerError, "Query failed", err)
		}
		var position string
		if result {
			position = "Inside"
		} else {
			position = "Outside"
		}
		responseBodyInfo := PolyResponse{geom, point, position}
		responseBody, err := json.Marshal(responseBodyInfo)
		if err != nil {
			c.Logger.Println("PolyResponse Marshal failed", err)
			c.WriteErrorResponse(w, http.StatusInternalServerError, "Could not marshal response", err)
			return
		}
		c.WriteResponse(w, http.StatusOK, responseBody)
	}
}

func (c *PolyController) DetermineGeogMembershipFromID() func(w http.ResponseWriter, r *http.Request) {

	type PolyResponse struct {
		Geom    *model.PolyGeometry `json:"geom"`
		Point    *model.PointGeometry  `json:"point"`
		Position string              `json:"position"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		idParams := mux.Vars(r)
		id, err := strconv.ParseInt(idParams["id"], 10, 0)
		if err != nil {
			c.WriteErrorResponse(w, http.StatusNotFound, "Invalid Path", err)
			return
		}
		intID := int(id)

		body, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			c.Logger.Println("Unprocessable request body", err)
			c.WriteErrorResponse(w, http.StatusInternalServerError, "Could not read body", err)
			return
		}

		var params IncomingFindClosestRequest
		err = json.Unmarshal(body, &params)
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
		pointJSON, err := json.Marshal(point)
		if err != nil {
			c.Logger.Println("Failed to Marshal pointJSON object")
			c.WriteErrorResponse(w, http.StatusInternalServerError, "Could not marshal geomJSON", err)
			return
		}
		pointString := string(pointJSON)
		queriedPolygon, err := c.Repository.GetPolygonFromID(intID)
		if err != nil {
			c.Logger.Println("Failed to retrieve polygon from given ID")
			c.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve polygon from given ID", err)
			return
		}

		result, err := c.Repository.Intersects(queriedPolygon, pointString)
		if err != nil {
			c.Logger.Println("DB Intersects Query failed")
			c.WriteErrorResponse(w, http.StatusInternalServerError, "Intersects Query failed", err)
			return
		}

		var position string
		if result {
			position = "Inside"
		} else {
			position = "Outside"
		}
		var resultGeom model.PolyGeometry
		err = json.Unmarshal([]byte(queriedPolygon), &resultGeom)
		if err != nil {
			c.Logger.Println("Failed to unmarshal response into polygon object")
			c.WriteErrorResponse(w, http.StatusInternalServerError, "Could not unmarshal geomJSON", err)
			return
		}
		responseBodyInfo := PolyResponse{&resultGeom, point, position}
		responseBody, err := json.Marshal(responseBodyInfo)
		if err != nil {
			c.Logger.Println("PolyResponse Marshal failed", err)
			c.WriteErrorResponse(w, http.StatusInternalServerError, "Could not marshal response", err)
			return
		}
		c.WriteResponse(w, http.StatusOK, responseBody)

	}
}

func (c *PolyController) FindMostProbableStore() func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			c.Logger.Println("Unprocessable request body", err)
			c.WriteErrorResponse(w, http.StatusInternalServerError, "Could not read body", err)
			return
		}

		var params IncomingFindClosestRequest
		err = json.Unmarshal(body, &params)
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

		result, err := c.Repository.FindClosest(params.StoreID, params.Point.Coordinates[0], params.Point.Coordinates[1])
		if err != nil {
			c.Logger.Println("DB Query failed")
			c.WriteErrorResponse(w, http.StatusInternalServerError, "DB Query failed", err)
			return
		}
		if result == (repository.LocationRow{}) {
			response := helpers.InsertResponse{"Unable to determine exact store"}
			responseBody, err := json.Marshal(response)
			if err != nil {
				c.Logger.Println("Message Marshal failed", err)
				c.WriteErrorResponse(w, http.StatusInternalServerError, "Could not marshal response", err)
				return
			}
			c.WriteResponse(w, http.StatusMultipleChoices, responseBody)
			return
		}
		responseBody, err := json.Marshal(result)
		if err != nil {
			c.Logger.Println("LocationRow Marshal failed", err)
			c.WriteErrorResponse(w, http.StatusInternalServerError, "Could not marshal response", err)
			return
		}
		c.WriteResponse(w, http.StatusOK, responseBody)
	}
}

func (c*PolyController) FindEnclosingPolygon() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			c.Logger.Println("Unprocessable request body", err)
			c.WriteErrorResponse(w, http.StatusInternalServerError, "Could not read body", err)
			return
		}

		var params IncomingFindEnclosingRequest
		err = json.Unmarshal(body, &params)
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

		result, err := c.Repository.FindEnclosingPolygon(params.Point.Coordinates[0], params.Point.Coordinates[1], params.StoreID, params.MetroID, params.ZoneID)
		if err != nil {
			c.Logger.Println("DB Query failed")
			c.WriteErrorResponse(w, http.StatusInternalServerError, "DB Query failed", err)
			return
		}
		if result == (repository.LocationRow{}) {
			response := helpers.InsertResponse{"Unable to determine exact store"}
			responseBody, err := json.Marshal(response)
			if err != nil {
				c.Logger.Println("Message Marshal failed", err)
				c.WriteErrorResponse(w, http.StatusInternalServerError, "Could not marshal response", err)
				return
			}
			c.WriteResponse(w, http.StatusMultipleChoices, responseBody)
			return
		}
		responseBody, err := json.Marshal(result)
		if err != nil {
			c.Logger.Println("LocationRow Marshal failed", err)
			c.WriteErrorResponse(w, http.StatusInternalServerError, "Could not marshal response", err)
			return
		}
		c.WriteResponse(w, http.StatusOK, responseBody)
	}
}

func (c *PolyController) Echo() func(w http.ResponseWriter, r *http.Request) {
	type IncomingMessage struct {
		Geom *model.PolyGeometry `json:"geom" validate:"required"`
		ID int `json:"id"`
	}
	return func(w http.ResponseWriter, r *http.Request) {

		body, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			c.Logger.Println("Unprocessable request body", err)
			c.WriteErrorResponse(w, http.StatusInternalServerError, "Could not read body", err)
			return
		}

		var params IncomingMessage
		err = json.Unmarshal(body, &params)
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

		id := params.ID
		geom := params.Geom

		responseBodyInfo := IncomingMessage{geom, id}
		responseBody, err := json.Marshal(responseBodyInfo)
		if err != nil {
			c.Logger.Println("PolyResponse Marshal failed", err)
			c.WriteErrorResponse(w, http.StatusInternalServerError, "Could not marshal response", err)
			return
		}
		c.WriteResponse(w, http.StatusOK, responseBody)
	}
}

func (c *PolyController) FeatureQuery() func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			c.Logger.Println("Unprocessable request body", err)
			c.WriteErrorResponse(w, http.StatusInternalServerError, "Could not read body", err)
			return
		}

		var params repository.LocationQuery
		err = json.Unmarshal(body, &params)
		if err != nil {
			c.Logger.Println("Failed to unmarshal LocationQuery", err)
			c.WriteErrorResponse(w, http.StatusInternalServerError, "Could not unmarshal input", err)
			return
		}

		if err != nil {
			c.Logger.Println("Unprocessable Request Body", err)
			c.WriteErrorResponse(w, http.StatusUnprocessableEntity, "Invalid Request Body", err)
			return
		}

		locationList, err := c.Repository.QueryDatabase(params)
		if err != nil {
			c.Logger.Println("Database Query Failed", err)
			c.WriteErrorResponse(w, http.StatusInternalServerError, "Query Failed", err)
		}
		c.Logger.Println(locationList)
		var featureList []interface{}
		featureList = helpers.ListToGeoJSONPointFeatures(locationList, c.Logger)

		responseBody, err := json.Marshal(featureList)
		if err != nil {
			c.Logger.Println("PolyLocationResponseCleaned Marshal failed", err)
			c.WriteErrorResponse(w, http.StatusInternalServerError, "Could not marshal response", err)
			return
		}
		c.WriteResponse(w, http.StatusOK, responseBody)
	}
}

func (c *PolyController) FindPolyLocationFromID() func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		idParams := mux.Vars(r)
		id, err := strconv.ParseInt(idParams["id"], 10, 0)
		if err != nil {
			c.WriteErrorResponse(w, http.StatusNotFound, "Invalid Path", err)
			return
		}
		intID := int(id)

		locationList, err := c.Repository.GetPolyLocationFromID(intID)
		if err != nil {
			c.Logger.Println("Database Query Failed", err)
			c.WriteErrorResponse(w, http.StatusInternalServerError, "Query Failed", err)
		}
		if len(locationList) == 0 {
			c.Logger.Println("No matching records found")
			c.WriteErrorResponse(w, http.StatusNoContent, "No matching record", nil)
		}
		var feature []interface{}
		feature = helpers.ListToGeoJSONPointFeatures(locationList, c.Logger)

		responseBody, err := json.Marshal(feature)
		if err != nil {
			c.Logger.Println("PolyLocationResponseCleaned Marshal failed", err)
			c.WriteErrorResponse(w, http.StatusInternalServerError, "Could not marshal response", err)
			return
		}
		c.WriteResponse(w, http.StatusOK, responseBody)

	}
}

