package helpers

import (
	"github.com/geofence/internal/json"
	"github.com/geofence/internal/model"
	"github.com/geofence/internal/repository"
	"log"
)

func AsGeoJSONPolyFeature(polyLocation repository.PolyLocationResponseCleaned, logger log.Logger) (repository.GeoJSONPolyFeature, error) {
	featureProperties := repository.FeatureProperties{
		ID: polyLocation.ID,
		Name:    polyLocation.Name,
		Street1: polyLocation.Street1,
		Zip:     polyLocation.Zip,
		City:    polyLocation.City,
		State:   polyLocation.State,
		MetroID: polyLocation.MetroID,
		ZoneID: polyLocation.ZoneID,
		StoreID: polyLocation.StoreID,
		County: polyLocation.County,
		OpeningHour: polyLocation.OpeningHour,
		Longitude: polyLocation.Longitude,
		Latitude: polyLocation.Latitude,
		Polygon: polyLocation.Polygon,
	}

	var geometry model.PolyGeometry
	err := json.Unmarshal([]byte(polyLocation.Polygon), &geometry)
	if err != nil {
		log.Println(err)
		return repository.GeoJSONPolyFeature{}, err
	}
	return repository.GeoJSONPolyFeature{
		Type: "Feature",
		Properties: featureProperties,
		Geometry: geometry,
	}, nil
}

func AsGeoJSONPointFeature(polyLocation repository.PolyLocationResponseCleaned) (repository.GeoJSONPointFeature) {
	featureProperties := repository.FeatureProperties{
		ID: polyLocation.ID,
		Name:    polyLocation.Name,
		Street1: polyLocation.Street1,
		Zip:     polyLocation.Zip,
		City:    polyLocation.City,
		State:   polyLocation.State,
		MetroID: polyLocation.MetroID,
		ZoneID: polyLocation.ZoneID,
		StoreID: polyLocation.StoreID,
		County: polyLocation.County,
		OpeningHour: polyLocation.OpeningHour,
		Longitude: polyLocation.Longitude,
		Latitude: polyLocation.Latitude,
		Polygon: polyLocation.Polygon,
	}
	geometry := model.PointGeometry{Type: "Point", Coordinates: [2]float64{polyLocation.Longitude, polyLocation.Latitude}}
	return repository.GeoJSONPointFeature{
		Type: "Feature",
		Properties: featureProperties,
		Geometry: geometry,
	}
}

func ListToGeoJSONFeatures(responses []repository.PolyLocationResponseCleaned, logger log.Logger) []interface{} {
	results := []interface{}{}
	for _, val := range responses {
		if val.Polygon == "" {
			results = append(results, AsGeoJSONPointFeature(val))
		} else {
			polyFeature, err := AsGeoJSONPolyFeature(val, logger)
			if err != nil {
				continue
			}
			results = append(results, polyFeature)
		}
	}
	return results
}

func ListToGeoJSONPointFeatures(responses []repository.PolyLocationResponseCleaned, logger log.Logger) []interface{} {
	results := []interface{}{}
	for _, val := range responses {
		pointFeature := AsGeoJSONPointFeature(val)
		results = append(results, pointFeature)
	}
	return results
}
