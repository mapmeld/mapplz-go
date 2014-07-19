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

func (mp *MapPLZ) Add_LatLng(latlng []float64) MapItem {
  mip := NewMapItemPoint(latlng[0], latlng[1])
	mp.MapItems = append(mp.MapItems, mip)
	return mip
}

func (mp *MapPLZ) Add_LngLat(lnglat []float64) MapItem {
	return mp.Add_Lat_Lng(lnglat[1], lnglat[0])
}

func (mp *MapPLZ) Add_Lat_Lng(lat float64, lng float64) MapItem {
	mip := NewMapItemPoint(lat, lng)
  mp.MapItems = append(mp.MapItems, mip)
	return mip
}

func (mp *MapPLZ) Add_Lng_Lat(lng float64, lat float64) MapItem {
	return mp.Add_Lat_Lng(lat, lng)
}

func (mp *MapPLZ) Add_LatLngPath(path [][]float64) MapItem {
	ml := NewMapItemLine(path)
	mp.MapItems = append(mp.MapItems, ml)
	return ml
}

func (mp *MapPLZ) Add_LngLatPath(lnglat_path [][]float64) MapItem {
	for i := 0; i < len(lnglat_path); i++ {
		lat := lnglat_path[i][1]
		lng := lnglat_path[i][0]
		lnglat_path[i][0] = lat
		lnglat_path[i][1] = lng
	}
	return mp.Add_LatLngPath(lnglat_path)
}

func (mp *MapPLZ) Add_LatLngPoly(path [][]float64) MapItem {
	ml := NewMapItemPoly(path)
	mp.MapItems = append(mp.MapItems, ml)
	return ml
}

func (mp *MapPLZ) Add_LngLatPoly(lnglat_path [][]float64) MapItem {
	for i := 0; i < len(lnglat_path); i++ {
		lat := lnglat_path[i][1]
		lng := lnglat_path[i][0]
		lnglat_path[i][0] = lat
		lnglat_path[i][1] = lng
	}
	return mp.Add_LatLngPoly(lnglat_path)
}


func (mp *MapPLZ) Add_Geojson_Collection(geojson []byte) GeojsonFeatureCollection {
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

func (mp *MapPLZ) Add_Geojson_Feature(geojson []byte) MapItem {
	// GeoJSON parsing based on http://stackoverflow.com/a/15728702
	var geojsonData GeojsonFeature
	err := json.Unmarshal(geojson, &geojsonData)
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
		panic("Unsupported type")
	}

	if err != nil {
		panic("Failed to parse JSON string")
	}

	return mip
}
