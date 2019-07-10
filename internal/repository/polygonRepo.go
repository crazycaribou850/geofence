package repository

import (
	"github.com/geofence/internal/model"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type PolygonPostgresRepository struct {
	DB sqlx.DB
}

func NewPolygonRepository(db sqlx.DB) *PolygonPostgresRepository {
	return &PolygonPostgresRepository{
		DB: db,
	}
}

// Takes in 2 strings which correctly represent geometry objects and returns if they intersect or not.
func (c *PolygonPostgresRepository) Intersects(item1 string, item2 string) (bool, error) {
	querySQL := `SELECT ST_Intersects(ST_GeomFromGeoJSON(:item1), ST_GeomFromGeoJSON(:item2))`
	intersectsRow := IntersectsRow{item1, item2}
	records, err := c.DB.NamedQuery(querySQL, intersectsRow)
	if err != nil {
		return false, err
	}
	defer records.Close()

	var result bool
	if records.Next() {
		err = records.Scan(&result)
	}
	return result, nil
}

func (c *PolygonPostgresRepository) InsertLocation(locationRequest LocationRow) (error) {
	insertSQL := `INSERT INTO store_locations (
		id,
		name,
		created_at,
		updated_at,
		street1,
		zip,
		city,
		state,
		metro_id,
		longitude,
		latitude,
		street2,
        zone_id,
		store_id,
		county,
		deleted_at,
		opening_hour,
		closing_hour,
		store_number,
		store_group,
		active,
		allows_pickup,
		is_envoy_only,
		service_area_id,
		sells_alcohol,
		tax_exempt
	)
	VALUES (
		:id,
		:name,
		:created_at,
		:updated_at,
		:street1,
		:zip,
		:city,
		:state,
		:metro_id,
		:longitude,
		:latitude,
		:street2,
        :zone_id,
		:store_id,
		:county,
		:deleted_at,
		:opening_hour,
		:closing_hour,
		:store_number,
		:store_group,
		:active,
		:allows_pickup,
		:is_envoy_only,
		:service_area_id,
		:sells_alcohol,
		:tax_exempt
	)
	`
	transaction, err := c.DB.Beginx()
	if err != nil {
		return err
	}
	rollback := false

	defer func() {
		if rollback {
			transaction.Rollback()
		} else {
			_ = transaction.Commit()
		}
	}()

	_, err = transaction.NamedExec(insertSQL, locationRequest)
	if err != nil {
		rollback = true
		return err
	}

	return nil
}

func (c *PolygonPostgresRepository) InsertPolygon(polygonID int, polygonObject model.PolyGeometry) (error) {
	insertSQL := `INSERT INTO store_polygons (
		id,
		polygon
	)
	VALUES (
		:id,
		ST_GeomFromGeoJSON(:polygon)
	)
	`

	row, err := toPolygonRow(polygonID, polygonObject)
	if err != nil {
		return err
	}

	transaction, err := c.DB.Beginx()
	if err != nil {
		return err
	}
	rollback := false

	defer func() {
		if rollback {
			transaction.Rollback()
		} else {
			_ = transaction.Commit()
		}
	}()

	_, err = transaction.NamedExec(insertSQL, row)
	if err != nil {
		rollback = true
		return err
	}

	return nil
}

func (c *PolygonPostgresRepository) Insert(insertRequest PolyLocationRow) (error) {
	insertSQL := `INSERT INTO store_locations (
		id,
		name,
		created_at,
		updated_at,
		street1,
		zip,
		city,
		state,
		metro_id,
		longitude,
		latitude,
		street2,
        zone_id,
		store_id,
		county,
		deleted_at,
		opening_hour,
		closing_hour,
		store_number,
		store_group,
		active,
		allows_pickup,
		is_envoy_only,
		service_area_id,
		sells_alcohol,
		tax_exempt
	)
	VALUES (
		:id,
		:name,
		:created_at,
		:updated_at,
		:street1,
		:zip,
		:city,
		:state,
		:metro_id,
		:longitude,
		:latitude,
		:street2,
        :zone_id,
		:store_id,
		:county,
		:deleted_at,
		:opening_hour,
		:closing_hour,
		:store_number,
		:store_group,
		:active,
		:allows_pickup,
		:is_envoy_only,
		:service_area_id,
		:sells_alcohol,
		:tax_exempt
	)
	`
	transaction, err := c.DB.Beginx()
	if err != nil {
		return err
	}
	rollback := false

	defer func() {
		if rollback {
			transaction.Rollback()
		} else {
			_ = transaction.Commit()
		}
	}()
	storeLocation := LocationRow{
		ID: insertRequest.ID,
		Name: insertRequest.Name,
		CreatedAt: insertRequest.CreatedAt,
		UpdatedAt: insertRequest.UpdatedAt,
		Street1: insertRequest.Street1,
		Zip: insertRequest.Zip,
		City: insertRequest.City,
		State: insertRequest.State,
		MetroID: insertRequest.MetroID,
		Longitude: insertRequest.Longitude,
		Latitude: insertRequest.Latitude,
		Street2: insertRequest.Street2,
		ZoneID: insertRequest.ZoneID,
		StoreID: insertRequest.StoreID,
		County: insertRequest.County,
		DeletedAt: insertRequest.DeletedAt,
		OpeningHour: insertRequest.OpeningHour,
		ClosingHour: insertRequest.ClosingHour,
		StoreNumber: insertRequest.StoreNumber,
		StoreGroup: insertRequest.StoreGroup,
		Active: insertRequest.Active,
		AllowsPickup: insertRequest.AllowsPickup,
		IsEnvoyOnly: insertRequest.IsEnvoyOnly,
		ServiceAreaId: insertRequest.ServiceAreaId,
		SellsAlcohol: insertRequest.SellsAlcohol,
		TaxExempt: insertRequest.TaxExempt,
	}
	_, err = transaction.NamedExec(insertSQL, storeLocation)
	if err != nil {
		rollback = true
		return err
	}

	insertSQL = `INSERT INTO store_polygons (
		id,
		polygon
	)
	VALUES (
		:id,
		ST_GeomFromGeoJSON(:polygon)
	)
	`

	row, err := toPolygonRow(insertRequest.ID, insertRequest.Polygon)
	if err != nil {
		rollback = true
		return err
	}
	_, err = transaction.NamedExec(insertSQL, row)
	if err != nil {
		rollback = true
		return err
	}

	return nil
}

func (c *PolygonPostgresRepository) GetAll() ([]PolyLocationResponseCleaned, error) {
	querySQL := `SELECT * FROM store_polygons NATURAL JOIN store_locations`

	records, err := c.DB.Queryx(querySQL)
	if err != nil {
		return []PolyLocationResponseCleaned{}, err
	}
	defer records.Close()

	var result PolyLocationResponse
	var results []PolyLocationResponse
	for records.Next() {
		err = records.StructScan(&result)
		if err != nil {
			return []PolyLocationResponseCleaned{}, err
		}
		results = append(results, result)
	}

	return PLResponseArrayToRegularTypes(results), nil
}

func (c *PolygonPostgresRepository) GetPolygonFromID(id int) (string, error) {
	querySQL := `SELECT ST_AsGeoJSON(polygon) FROM store_polygons WHERE id = :id`
	IDRow := IDRow{ID: id}
	records, err := c.DB.NamedQuery(querySQL, IDRow)
	if err != nil {
		return "", err
	}
	defer records.Close()

	var result string
	if records.Next() {
		err = records.Scan(&result)
		if err != nil {
			return "", err
		}
	}
	if result == "" {
		return result, errors.New("No polygon with that ID found")
	}
	return result, nil
}

func (c *PolygonPostgresRepository) GetPolygonFromStoreID(storeID string) ([]PolyLocationResponseCleaned, error) {
	querySQL := `SELECT * FROM store_polygons as sp, store_locations as sl WHERE sl.store_id = :store_id AND sl.id = sp.id`
	SIDRow := SIDRow{StoreID: storeID}
	records, err := c.DB.NamedQuery(querySQL, SIDRow)
	if err != nil {
		return []PolyLocationResponseCleaned{}, err
	}
	defer records.Close()

	var results []PolyLocationResponse
	var result PolyLocationResponse
	for records.Next() {
		err = records.StructScan(&result)
		if err != nil {
			return []PolyLocationResponseCleaned{}, err
		}
		results = append(results, result)
	}

	return PLResponseArrayToRegularTypes(results), nil
}





