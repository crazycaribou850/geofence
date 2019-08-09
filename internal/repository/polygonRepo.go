package repository

import (
	"github.com/geofence/internal/model"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"strconv"
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
	ON CONFLICT (id) DO UPDATE SET polygon = ST_GeomFromGeoJSON(:polygon)
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

func (c *PolygonPostgresRepository) GetAll() ([]PolyLocationResponseCleaned, error) {
	querySQL := `SELECT * FROM store_polygons NATURAL JOIN store_locations`
	var results []PolyLocationResponse
	err := c.DB.Select(&results, querySQL)
	if err != nil {
		return []PolyLocationResponseCleaned{}, err
	}
	return PLResponseArrayToRegularTypes(results), nil
}

func (c *PolygonPostgresRepository) GetPolygonFromID(id int) (string, error) {
	querySQL := `SELECT ST_AsGeoJSON(polygon) FROM store_polygons WHERE id = $1`
	var result string
	err := c.DB.Select(&result, querySQL, id)
	if err != nil {
		return "", err
	}
	if result == "" {
		return result, errors.New("No polygon with that ID found")
	}
	return result, nil
}

func (c *PolygonPostgresRepository) GetPolyLocationFromID(id int) ([]PolyLocationResponseCleaned, error) {
	querySQL := `SELECT sl.*, ST_AsGeoJSON(sp.polygon) as polygon FROM store_locations as sl LEFT JOIN store_polygons as sp ON (sl.id = sp.id) WHERE sl.id = $1`
	var result []PolyLocationResponse
	err := c.DB.Select(&result, querySQL, id)
	if err != nil {
		return []PolyLocationResponseCleaned{}, err
	}
	if len(result) == 0 {
		return []PolyLocationResponseCleaned{}, nil
	}
	return PLResponseArrayToRegularTypes(result), nil
}

func (c *PolygonPostgresRepository) QueryDatabase(data LocationQuery) ([]PolyLocationResponseCleaned, error) {
	var appendedCount int
	var storeIDclause string
	var metroIDclause string
	var zoneIDclause string
	var cityClause string
	var stateClause string
	if (data.ID != 0) {
		return c.GetPolyLocationFromID(data.ID)
	}
	if (data.StoreID != 0) {
		storeIDclause = `sl.store_id = ` + strconv.Itoa(data.StoreID)
	}
	if (data.MetroID != 0) {
		metroIDclause = `sl.metro_id = ` + strconv.Itoa(data.MetroID)
	}
	if (data.ZoneID != 0) {
		zoneIDclause = `sl.zone_id = ` + strconv.Itoa(data.ZoneID)
	}
	if (data.City != "") {
		cityClause = `sl.city = '` + data.City + `'`
	}
	if (data.State != "") {
		stateClause = `sl.state = '` + data.State + `'`
	}
	baseQuery := `SELECT sl.*, ST_AsGeoJSON(sp.polygon) as polygon FROM store_locations as sl LEFT JOIN store_polygons as sp ON (sl.id = sp.id)`
	baseQuery = appendClause(baseQuery, storeIDclause, &appendedCount)
	baseQuery = appendClause(baseQuery, metroIDclause, &appendedCount)
	baseQuery = appendClause(baseQuery, zoneIDclause, &appendedCount)
	baseQuery = appendClause(baseQuery, cityClause, &appendedCount)
	baseQuery = appendClause(baseQuery, stateClause, &appendedCount)
	baseQuery = baseQuery + ` ORDER BY sl.id`
	var results []PolyLocationResponse

	err := c.DB.Select(&results, baseQuery)
	if err != nil {
		return []PolyLocationResponseCleaned{}, err
	}

	return PLResponseArrayToRegularTypes(results), nil
}

func appendClause(baseQuery, clause string, appended *int) string {
	if clause == "" {
		return baseQuery
	}
	if *appended == 0 {
		*appended += 1
		return baseQuery + ` WHERE ` + clause
	} else {
		*appended += 1
		return baseQuery + ` AND ` + clause
	}
}
func (c*PolygonPostgresRepository) FindClosest(store_id int, lat, long float64) (LocationRow, error) {
	querySQL := `WITH candidates (id, distance) AS (SELECT id, ST_Distance(ST_MakePoint(latitude, longitude), ST_MakePoint($2, $3)) as distance FROM store_locations 
					WHERE ST_DWithin(ST_MakePoint(latitude, longitude), ST_MakePoint($2, $3), 1) AND active=True AND store_id= $1)
					SELECT store_locations.* FROM candidates, store_locations
					WHERE store_locations.id = candidates.id AND candidates.distance in (SELECT MIN(candidates.distance) FROM candidates)`
	var results []LocationRowNull
	err := c.DB.Select(&results, querySQL, store_id, lat, long)
	if err != nil {
		return LocationRow{}, err
	}
	if len(results) == 0 {
		return LocationRow{}, nil
	}
	if len(results) == 1 {
		return LocationToRegularTypes(results[0]), nil
	} else {
		result, err := c.checkPolygons(results, lat, long)
		if err != nil {
			return LocationRow{}, err
		}
		return LocationToRegularTypes(result), nil
	}
}

func (c*PolygonPostgresRepository) checkPolygons(rows []LocationRowNull, lat, long float64) (LocationRowNull, error) {
	var indices []int
	for _, row := range rows {
		indices = append(indices, row.ID)
	}
	proceed, err := c.checkAllExistsInPolygonTable(indices)
	if err != nil {
		return LocationRowNull{}, err
	}
	if proceed == false {
		return LocationRowNull{}, nil
	}
	params := map[string]interface{}{
		"lat": lat,
		"long": long,
		"ids": indices,

	}
	querySQL := `SELECT sl.* FROM store_locations sl, store_polygons sp WHERE sl.id IN (:ids) AND sp.id = sl.id AND ST_Intersects(sp.polygon, ST_MakePoint(:lat, :long))`
	querySQL, args, err := sqlx.Named(querySQL, params)
	if err != nil {
		return LocationRowNull{}, err
	}
	querySQL, args, err = sqlx.In(querySQL, args)
	if err != nil {
		return LocationRowNull{}, err
	}
	querySQL = c.DB.Rebind(querySQL)
	var results []LocationRowNull
	err = c.DB.Select(&results, querySQL)
	if err != nil {
		return LocationRowNull{}, nil
	}
	if len(results) != 1 {
		return LocationRowNull{}, nil
	}
	return results[0], nil
}

func (c*PolygonPostgresRepository) checkAllExistsInPolygonTable(indices []int) (bool, error) {
	querySQL := `SELECT * FROM store_polygons WHERE id IN (?)`
	query, args, err := sqlx.In(querySQL, indices)
	if err != nil {
		return false, err
	}
	c.DB.Rebind(query)
	var records []PolygonRow
	err = c.DB.Select(records, query, args)
	if err != nil {
		return false, err
	}
	if len(records) == len(indices) {
		return true, nil
	} else {
		return false, nil
	}
}

// Assumes that all polygons have been drawn for
func (c*PolygonPostgresRepository) FindEnclosingPolygon(lat, long float64, storeID, metroID, zoneID int) (LocationRow, error) {
	querySQL := `SELECT sl.* FROM store_locations sl, store_polygons sp 
				WHERE sl.id=sp.id AND sl.store_id=$3 AND sl.metro_id=$4 AND sl.zone_id=$5 AND ST_Intersects(ST_MakePoint($1, $2), sp.polygon)`
	var results []LocationRow
	err := c.DB.Select(results, querySQL, lat, long, storeID, metroID, zoneID)
	if err != nil {
		return LocationRow{}, err
	}
	if len(results) != 1 {
		return LocationRow{}, nil
	} else {
		return results[0], nil
	}
}




