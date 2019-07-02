package repository
//
//import (
//	"context"
//	"github.com/geofence/internal/model"
//	"github.com/newshipt/shipt-tofu/database/sqlt"
//)
//
//type PolygonRepository interface {
//	Upsert(ctx context.Context, polygonRequest model.PolygonRequest) (model.Polygon, error)
//}
//
//type PolygonPostgresRepository struct {
//	PostgresClientDriver sqlt.SQL
//}
//
//func NewPolygonRepository(postgresClientDriver sqlt.SQL) *PolygonPostgresRepository {
//	return &PolygonPostgresRepository{
//		PostgresClientDriver: postgresClientDriver,
//	}
//}
//
//func (c *PolygonPostgresRepository) Upsert(ctx context.Context, polygonRequest model.PolygonRequest) (model.Polygon, error) {
//	insertSQL := `INSERT INTO polygons (
//		id,
//		name,
//		store_code,
//		metro_id,
//		zone_id,
//		polygon,
//		point
//	)
//	VALUES (
//		:id ,
//		:name,
//		:store_code,
//		:metro_id,
//		:zone_id,
//		:polygon,
//		:point
//	)
//	ON CONFLICT (id) DO UPDATE
//	SET name = EXCLUDED.name,
//		store_code = EXCLUDED.store_code,
//		metro_id = EXCLUDED.metro_id,
//		zone_id = EXCLUDED.zone_id,
//		polygon = EXCLUDED.polygon,
//		point = ST_Centroid(polygon)
//	RETURNING id, name, metro_id, zone_id, polygon, opoint.
//		`
//}
//
//var result Polygon
//row := toPolygonRow(polygonRequest)
//
//transaction, err := c.PostgresClient.Driver.Begin(ctx)
//if err != nil {
//	return model.Polygon{}, err
//}
//
//
//
