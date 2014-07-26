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
	Features []string
}

func (mp *MapPLZ) ToGeoJson() string {
	if mp.Database != nil {
		// fetch all items from the DB
		mp.MapItems = mp.Query("")
	}
	var features = []string{}
	for i := 0; i < len(mp.MapItems); i++ {
		features = append(features, mp.MapItems[i].ToGeoJson())
	}
	return `{"type":"FeatureCollection","features":[` + strings.Join(features, ",") + `]}`
}

func (mp *MapPLZ) Add_Geojson_Feature(geojson string) MapItem {
	mip := ConvertGeojsonFeature(geojson, mp.Database)
	mp.MapItems = append(mp.MapItems, mip)

	return mip
}

func (mp *MapPLZ) Add_Geojson_Collection(geojson string) []MapItem {
	// GeoJSON parsing based on http://stackoverflow.com/a/15728702
	var geojsonData GeojsonFeatureCollection
	var featureList []MapItem
	json.Unmarshal([]byte(geojson), &geojsonData)

	for i := range geojsonData.Features {
		t := geojsonData.Features[i]
		featureList = append(featureList, ConvertGeojsonFeature(t, mp.Database))
	}

	return featureList
}
