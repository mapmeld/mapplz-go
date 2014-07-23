package mapplz

import (
	"encoding/json"
	"strings"
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

func (mp *MapPLZ) ToGeoJson() string {
	var features = []string{}
	for i := 0; i < len(mp.MapItems); i++ {
		features = append(features, mp.MapItems[i].ToGeoJson())
	}
	return `{"type":"FeatureCollection","features":[` + strings.Join(features, ",") + `]}`
}

func (mp *MapPLZ) Add_Geojson_Feature(geojson string) MapItem {
	mip := ConvertGeojsonFeature(geojson)
	mp.MapItems = append(mp.MapItems, mip)

	return mip
}

func (mp *MapPLZ) Add_Geojson_Collection(geojson string) GeojsonFeatureCollection {
	// GeoJSON parsing based on http://stackoverflow.com/a/15728702
	var geojsonData GeojsonFeatureCollection
	err := json.Unmarshal([]byte(geojson), &geojsonData)

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
