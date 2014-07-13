package mapplz

import (
	"encoding/json"
	"github.com/kellydunn/golang-geo"
	// "github.com/kpawlik/geojson"
)

type GeojsonPoint struct {
	Coordinates []float64
}

type GeojsonLine struct {
	Coordinates [][]float64
}

type GeojsonPolygon struct {
	Coordinates [][][]float64
}

type GeojsonGeometry struct {
	Type        string
	Coordinates json.RawMessage
	Point       GeojsonPoint
	Line        GeojsonLine
	Polygon     GeojsonPolygon
}

type GeojsonFeature struct {
	Type       string
	Geometry   GeojsonGeometry
	Properties json.RawMessage
}

type GeojsonFeatureCollection struct {
	Type     string
	Features []GeojsonFeature
}

func Add_LatLng(latlng []float64) *geo.Point {
	pt := geo.NewPoint(latlng[0], latlng[1])
	return pt
}

func Add_LngLat(lnglat []float64) *geo.Point {
	return Add_Lat_Lng(lnglat[1], lnglat[0])
}

func Add_Lat_Lng(lat float64, lng float64) *geo.Point {
	pt := geo.NewPoint(lat, lng)
	return pt
}

func Add_Lng_Lat(lng float64, lat float64) *geo.Point {
	return Add_Lat_Lng(lat, lng)
}

func Add_Geojson_Collection_Str(geojson []byte) GeojsonFeatureCollection {
	// GeoJSON parsing based on http://stackoverflow.com/a/15728702
	var geojsonData GeojsonFeatureCollection
	err := json.Unmarshal(geojson, &geojsonData)

	for i := range geojsonData.Features {
		t := &geojsonData.Features[i]

		switch t.Geometry.Type {
		case "Point":
			err = json.Unmarshal(t.Geometry.Coordinates, &t.Geometry.Point.Coordinates)
		case "LineString":
			err = json.Unmarshal(t.Geometry.Coordinates, &t.Geometry.Line.Coordinates)
		case "Polygon":
			err = json.Unmarshal(t.Geometry.Coordinates, &t.Geometry.Polygon.Coordinates)
		default:
			panic("Unsupported type")
		}

		if err != nil {
			panic("Failed to parse JSON string")
		}

	}

	return geojsonData
}

func Add_Geojson_Feature_Str(geojson []byte) GeojsonFeature {
	// GeoJSON parsing based on http://stackoverflow.com/a/15728702
	var geojsonData GeojsonFeature
	err := json.Unmarshal(geojson, &geojsonData)

	switch geojsonData.Geometry.Type {
	case "Point":
		err = json.Unmarshal(geojsonData.Geometry.Coordinates, &geojsonData.Geometry.Point.Coordinates)
	case "LineString":
		err = json.Unmarshal(geojsonData.Geometry.Coordinates, &geojsonData.Geometry.Line.Coordinates)
	case "Polygon":
		err = json.Unmarshal(geojsonData.Geometry.Coordinates, &geojsonData.Geometry.Polygon.Coordinates)
	default:
		panic("Unsupported type")
	}

	if err != nil {
		panic("Failed to parse JSON string")
	}

	return geojsonData
}
