package mapplz

import (
	"encoding/json"
)

type MapPLZ struct {
  MapItems   []MapItem
}

func NewMapPLZ() MapPLZ {
  var mis = []MapItem{}
  return MapPLZ{MapItems: mis}
}

type MapItem interface {
  Type()       string
  Lat()        float64
  Lng()        float64
  Path()       [][][]float64
  Properties() json.RawMessage
  ToGeoJson()  string
}

func Add_LatLng(latlng []float64) MapItem {
  mip := NewMapItemPoint(latlng[0], latlng[1])
	return mip
}

func Add_LngLat(lnglat []float64) MapItem {
	return Add_Lat_Lng(lnglat[1], lnglat[0])
}

func Add_Lat_Lng(lat float64, lng float64) MapItem {
  mip := NewMapItemPoint(lat, lng)
  return mip
}

func (mp MapPLZ) Add_Lat_Lng(lat float64, lng float64) (MapPLZ, MapItem) {
	mip := NewMapItemPoint(lat, lng)
  mp.MapItems = append(mp.MapItems, mip)
	return mp, mip
}

func Add_Lng_Lat(lng float64, lat float64) MapItem {
	return Add_Lat_Lng(lat, lng)
}

func Add_Geojson_Collection(geojson []byte) GeojsonFeatureCollection {
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

func Add_Geojson_Feature(geojson []byte) GeojsonFeature {
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
