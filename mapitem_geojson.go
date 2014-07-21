package mapplz

import (
	"encoding/json"
	"fmt"
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
	// GeoJSON parsing based on http://stackoverflow.com/a/15728702
	var geojsonData GeojsonFeature
	err := json.Unmarshal([]byte(geojson), &geojsonData)
	var mip MapItem

	switch geojsonData.Geometry.Type {
	case "Point":
		err = json.Unmarshal(geojsonData.Geometry.Coordinates, &geojsonData.Geometry.Point.Coordinates)
		mip = mp.Add_LngLat(geojsonData.Geometry.Point.Coordinates)
	case "LineString":
		err = json.Unmarshal(geojsonData.Geometry.Coordinates, &geojsonData.Geometry.Line.Coordinates)
		mip = mp.Add_LngLatPath(geojsonData.Geometry.Line.Coordinates)
	case "Polygon":
		err = json.Unmarshal(geojsonData.Geometry.Coordinates, &geojsonData.Geometry.Polygon.Coordinates)
		mip = mp.Add_LngLatPoly(geojsonData.Geometry.Polygon.Coordinates[0])
	default:
		fmt.Printf("%s", geojson)
		panic("Unsupported type")
	}

	if err != nil {
		panic("Failed to parse JSON string")
	}

	mip.SetProperties(geojsonData.Properties)

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
