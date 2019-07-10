package repository

import (
	"database/sql"
	"github.com/geofence/internal/json"
	"github.com/geofence/internal/model"
	"github.com/lib/pq"
	"time"
)

type LocationRow struct {
	ID		int `db:"id" validate:"required"`
	Name	string `db:"name"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	Street1 string	`db:"street1"`
	Zip		string 	`db:"zip"`
	City	string 	`db:"city"`
	State	string	`db:"state"`
	MetroID int `db:"metro_id"`
	Longitude float64 `db:"longitude"`
	Latitude float64 `db:"latitude"`
	Street2 string `db:"street2"`
	ZoneID	int `db:"zone_id"`
	StoreID int `db:"store_id"`
	County string `db:"county"`
	DeletedAt time.Time `db:"deleted_at"`
	OpeningHour int `db:"opening_hour"`
	ClosingHour int `db:"closing_hour"`
	StoreNumber string `db:"store_number"`
	StoreGroup string `db:"store_group"`
	Active bool `db:"active"`
	AllowsPickup bool `db:"allows_pickup"`
	IsEnvoyOnly bool `db:"is_envoy_only"`
	ServiceAreaId int `db:"service_area_id"`
	SellsAlcohol bool `db:"sells_alcohol"`
	TaxExempt bool `db:"tax_exempt"`
}

type PolygonRow struct {
	ID	int 	`db:"id" validate:"required"`
	Polygon string 	`db:"polygon" validate:"required"`
}

type CoordinateRow struct {
	ID int `db:"id" validate:"required"`
	Point string `db:"point" validate:"required"`
}

type PolyLocationRow struct {
	ID		int `db:"id" validate:"required"`
	Name	string `db:"name"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	Street1 string	`db:"street1"`
	Zip		string 	`db:"zip"`
	City	string 	`db:"city"`
	State	string	`db:"state"`
	MetroID int `db:"metro_id"`
	Longitude float64 `db:"longitude"`
	Latitude float64 `db:"latitude"`
	Street2 string `db:"street2"`
	ZoneID	int `db:"zone_id"`
	StoreID int `db:"store_id"`
	County string `db:"county"`
	DeletedAt time.Time `db:"deleted_at"`
	OpeningHour int `db:"opening_hour"`
	ClosingHour int `db:"closing_hour"`
	StoreNumber string `db:"store_number"`
	StoreGroup string `db:"store_group"`
	Active bool `db:"active"`
	AllowsPickup bool `db:"allows_pickup"`
	IsEnvoyOnly bool `db:"is_envoy_only"`
	ServiceAreaId int `db:"service_area_id"`
	SellsAlcohol bool `db:"sells_alcohol"`
	TaxExempt bool `db:"tax_exempt"`
	Polygon model.PolyGeometry 	`json:"polygon" db:"polygon" validate:"required"`
}

type PolyLocationResponse struct {
	ID		int `db:"id" validate:"required"`
	Name	sql.NullString `db:"name"`
	CreatedAt pq.NullTime `db:"created_at"`
	UpdatedAt pq.NullTime `db:"updated_at"`
	Street1 sql.NullString	`db:"street1"`
	Zip		sql.NullString 	`db:"zip"`
	City	sql.NullString 	`db:"city"`
	State	sql.NullString	`db:"state"`
	MetroID sql.NullInt64 `db:"metro_id"`
	Longitude sql.NullFloat64 `db:"longitude"`
	Latitude sql.NullFloat64 `db:"latitude"`
	Street2 sql.NullString `db:"street2"`
	ZoneID	sql.NullInt64 `db:"zone_id"`
	StoreID sql.NullInt64 `db:"store_id"`
	County sql.NullString `db:"county"`
	DeletedAt pq.NullTime `db:"deleted_at"`
	OpeningHour sql.NullInt64 `db:"opening_hour"`
	ClosingHour sql.NullInt64 `db:"closing_hour"`
	StoreNumber sql.NullString `db:"store_number"`
	StoreGroup sql.NullString `db:"store_group"`
	Active sql.NullBool `db:"active"`
	AllowsPickup sql.NullBool `db:"allows_pickup"`
	IsEnvoyOnly sql.NullBool `db:"is_envoy_only"`
	ServiceAreaId sql.NullInt64 `db:"service_area_id"`
	SellsAlcohol sql.NullBool `db:"sells_alcohol"`
	TaxExempt sql.NullBool `db:"tax_exempt"`
	Polygon sql.NullString 	`json:"polygon" db:"polygon" validate:"required"`
}

type PolyLocationResponseCleaned struct {
	ID		int `db:"id" validate:"required"`
	Name	string `db:"name"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	Street1 string	`db:"street1"`
	Zip		string 	`db:"zip"`
	City	string 	`db:"city"`
	State	string	`db:"state"`
	MetroID int64 `db:"metro_id"`
	Longitude float64 `db:"longitude"`
	Latitude float64 `db:"latitude"`
	Street2 string `db:"street2"`
	ZoneID	int64 `db:"zone_id"`
	StoreID int64 `db:"store_id"`
	County string `db:"county"`
	DeletedAt time.Time `db:"deleted_at"`
	OpeningHour int64 `db:"opening_hour"`
	ClosingHour int64 `db:"closing_hour"`
	StoreNumber string `db:"store_number"`
	StoreGroup string `db:"store_group"`
	ServiceAreaId int64 `db:"service_area_id"`
	Polygon string 	`json:"polygon" db:"polygon" validate:"required"`
}

type IntersectsRow struct {
	Item1	string `db:"item1"`
	Item2	string `db:"item2"`
}

type IDRow struct {
	ID int `db:"id"`
}

type SIDRow struct {
	StoreID string `db:"store_id"`
}

func toPolygonRow(polygonID int, polygonObject model.PolyGeometry) (*PolygonRow, error) {
	polyGeom, err := json.Marshal(polygonObject)
	if err != nil {
		return nil, err
	}
	row := &PolygonRow{
		ID:			polygonID,
		Polygon:	string(polyGeom),
	}
	return row, nil
}

func PLResponseToRegularTypes(response PolyLocationResponse) (PolyLocationResponseCleaned) {
	return PolyLocationResponseCleaned{
		ID: response.ID,
		Name: response.Name.String,
		CreatedAt: response.CreatedAt.Time,
		UpdatedAt: response.UpdatedAt.Time,
		Street1: response.Street1.String,
		Zip: response.Zip.String,
		City: response.City.String,
		State: response.State.String,
		MetroID: response.MetroID.Int64,
		Longitude: response.Longitude.Float64,
		Latitude: response.Latitude.Float64,
		Street2: response.Street2.String,
		ZoneID: response.ZoneID.Int64,
		StoreID: response.StoreID.Int64,
		County: response.County.String,
		DeletedAt: response.DeletedAt.Time,
		OpeningHour: response.OpeningHour.Int64,
		ClosingHour: response.ClosingHour.Int64,
		StoreNumber: response.StoreNumber.String,
		StoreGroup: response.StoreGroup.String,
		ServiceAreaId: response.ServiceAreaId.Int64,
		Polygon: response.Polygon.String,
	}
}

func PLResponseArrayToRegularTypes(responseArray []PolyLocationResponse) ([]PolyLocationResponseCleaned) {
	var result []PolyLocationResponseCleaned
	for _, response := range responseArray {
		result = append(result, PLResponseToRegularTypes(response))
	}
	return result
}