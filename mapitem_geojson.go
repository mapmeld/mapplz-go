package mapplz

import (
	"encoding/json"
)

type geojsonPoint struct {
	Coordinates []float64
}

type geojsonLine struct {
	Coordinates [][]float64
}

type geojsonPolygon struct {
	Coordinates [][][]float64
}

type geojsonGeometry struct {
	Type        string
	Coordinates json.RawMessage
	Point       geojsonPoint
	Line        geojsonLine
	Polygon     geojsonPolygon
}

type GeojsonFeature struct {
	Type       string
	Geometry   geojsonGeometry
	Properties map[string]interface{}
}

type GeojsonFeatureCollection struct {
	Type     string
	Features []GeojsonFeature
}
