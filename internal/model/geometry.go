package model

type Polygon struct {
	Name	string `json:"name" validate:"required"`
	Polygon PolyGeometry `json:"polygon" validate:"required"`
}

type PolyGeometry struct {
	Type string `json:"type" validate:"required"`
	Coordinates [][][2]float64 `json:"coordinates" validate:"required"`
}

type PointGeometry struct {
	Type string `json:"type" validate:"required"`
	Coordinates [2]float64 `json:"coordinates" validate:"required"`
}