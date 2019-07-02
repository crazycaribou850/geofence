package model


type Polygon struct {
	ID      string `json:"id,omitempty"`
	Name	string `json:"string" validate:"required"`
	StoreCode string `json:"store_code" validate:"required"`
	MetroID string `json:"metro_id" validate:"required"`
	ZoneID	string `json:"zone_id" validate:"required"`
	Polygon PolyGeometry `json:"polygon" validate:"rquired"`
}

type PolyGeometry struct {
	Type string `json:"type" validate:"required"`
	Coordinates [][][2]float64 `json:"coordinates"`
}
