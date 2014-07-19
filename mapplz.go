package mapplz

import (
	"encoding/json"
)

type MapPLZ struct {
	MapItems []MapItem
}

func NewMapPLZ() MapPLZ {
	var mis = []MapItem{}
	return MapPLZ{MapItems: mis}
}

type MapItem interface {
	Type() string
	Lat() float64
	Lng() float64
	Path() [][][]float64
	SetProperties(map[string]interface{})
	Properties() map[string]interface{}
	SetJsonProperties(string)
	ToGeoJson() string
}

// lat, lng with variations

func (mp *MapPLZ) Add_Lat_Lng(lat float64, lng float64) MapItem {
	mip := NewMapItemPoint(lat, lng)
	mp.MapItems = append(mp.MapItems, mip)
	return mip
}

func (mp *MapPLZ) Add_Lat_Lng_Properties(lat float64, lng float64, props map[string]interface{}) MapItem {
	mip := NewMapItemPoint(lat, lng)
	mip.SetProperties(props)
	mp.MapItems = append(mp.MapItems, mip)
	return mip
}

func (mp *MapPLZ) Add_Lat_Lng_Json(lat float64, lng float64, props string) MapItem {
	var prop_map = map[string]interface{}{}
	json.Unmarshal([]byte(props), &prop_map)
	return mp.Add_Lat_Lng_Properties(lat, lng, prop_map)
}

// lng, lat with variations

func (mp *MapPLZ) Add_Lng_Lat(lng float64, lat float64) MapItem {
	return mp.Add_Lat_Lng(lat, lng)
}

func (mp *MapPLZ) Add_Lng_Lat_Properties(lng float64, lat float64, props map[string]interface{}) MapItem {
	return mp.Add_Lat_Lng_Properties(lat, lng, props)
}

func (mp *MapPLZ) Add_Lng_Lat_Json(lng float64, lat float64, props string) MapItem {
	return mp.Add_Lat_Lng_Json(lat, lng, props)
}

// [lat, lng] with variations

func (mp *MapPLZ) Add_LatLng(latlng []float64) MapItem {
	return mp.Add_Lat_Lng(latlng[0], latlng[1])
}

func (mp *MapPLZ) Add_LatLng_Properties(latlng []float64, props map[string]interface{}) MapItem {
	return mp.Add_Lat_Lng_Properties(latlng[0], latlng[1], props)
}

func (mp *MapPLZ) Add_LatLng_Json(latlng []float64, props string) MapItem {
	return mp.Add_Lat_Lng_Json(latlng[0], latlng[1], props)
}

func (mp *MapPLZ) Add_LatLngProperties(latlngprops []interface{}) MapItem {
	var prop_map = latlngprops[2].(map[string]interface{})
	return mp.Add_Lat_Lng_Properties(latlngprops[0].(float64), latlngprops[1].(float64), prop_map)
}

func (mp *MapPLZ) Add_LatLngJson(latlngprops []interface{}) MapItem {
	return mp.Add_Lat_Lng_Json(latlngprops[0].(float64), latlngprops[1].(float64), latlngprops[2].(string))
}

// [lng, lat] with variations

func (mp *MapPLZ) Add_LngLat(lnglat []float64) MapItem {
	return mp.Add_Lat_Lng(lnglat[1], lnglat[0])
}

func (mp *MapPLZ) Add_LngLat_Properties(lnglat []float64, props map[string]interface{}) MapItem {
	return mp.Add_Lat_Lng_Properties(lnglat[1], lnglat[0], props)
}

func (mp *MapPLZ) Add_LngLat_Json(lnglat []float64, props string) MapItem {
	return mp.Add_Lat_Lng_Json(lnglat[1], lnglat[0], props)
}

func (mp *MapPLZ) Add_LngLatProperties(lnglatprops []interface{}) MapItem {
	var prop_map = lnglatprops[2].(map[string]interface{})
	return mp.Add_Lat_Lng_Properties(lnglatprops[1].(float64), lnglatprops[0].(float64), prop_map)
}

func (mp *MapPLZ) Add_LngLatJson(lnglatprops []interface{}) MapItem {
	return mp.Add_Lat_Lng_Json(lnglatprops[1].(float64), lnglatprops[0].(float64), lnglatprops[2].(string))
}

// latlng path

func (mp *MapPLZ) Add_LatLngPath(path [][]float64) MapItem {
	ml := NewMapItemLine(path)
	mp.MapItems = append(mp.MapItems, ml)
	return ml
}

func (mp *MapPLZ) Add_LatLngPath_Properties(path [][]float64, props map[string]interface{}) MapItem {
	ml := NewMapItemLine(path)
	ml.SetProperties(props)
	mp.MapItems = append(mp.MapItems, ml)
	return ml
}

func (mp *MapPLZ) Add_LatLngPath_Json(path [][]float64, props string) MapItem {
	var prop_map = map[string]interface{}{}
	json.Unmarshal([]byte(props), &prop_map)
	return mp.Add_LatLngPath_Properties(path, prop_map)
}

// lnglat path

func (mp *MapPLZ) Add_LngLatPath(lnglat_path [][]float64) MapItem {
	for i := 0; i < len(lnglat_path); i++ {
		lat := lnglat_path[i][1]
		lng := lnglat_path[i][0]
		lnglat_path[i][0] = lat
		lnglat_path[i][1] = lng
	}
	return mp.Add_LatLngPath(lnglat_path)
}

func (mp *MapPLZ) Add_LngLatPath_Properties(path [][]float64, props map[string]interface{}) MapItem {
	ml := mp.Add_LngLatPath(path)
	ml.SetProperties(props)
	return ml
}

func (mp *MapPLZ) Add_LngLatPath_Json(path [][]float64, props string) MapItem {
	var prop_map = map[string]interface{}{}
	json.Unmarshal([]byte(props), &prop_map)
	return mp.Add_LngLatPath_Properties(path, prop_map)
}

// latlng poly

func (mp *MapPLZ) Add_LatLngPoly(path [][]float64) MapItem {
	ml := NewMapItemPoly(path)
	mp.MapItems = append(mp.MapItems, ml)
	return ml
}

func (mp *MapPLZ) Add_LatLngPoly_Properties(path [][]float64, props map[string]interface{}) MapItem {
	ml := NewMapItemPoly(path)
	ml.SetProperties(props)
	mp.MapItems = append(mp.MapItems, ml)
	return ml
}

func (mp *MapPLZ) Add_LatLngPoly_Json(path [][]float64, props string) MapItem {
	var prop_map = map[string]interface{}{}
	json.Unmarshal([]byte(props), &prop_map)
	return mp.Add_LatLngPoly_Properties(path, prop_map)
}

// lnglat poly

func (mp *MapPLZ) Add_LngLatPoly(lnglat_path [][]float64) MapItem {
	for i := 0; i < len(lnglat_path); i++ {
		lat := lnglat_path[i][1]
		lng := lnglat_path[i][0]
		lnglat_path[i][0] = lat
		lnglat_path[i][1] = lng
	}
	return mp.Add_LatLngPoly(lnglat_path)
}

func (mp *MapPLZ) Add_LngLatPoly_Properties(path [][]float64, props map[string]interface{}) MapItem {
	ml := mp.Add_LngLatPoly(path)
	ml.SetProperties(props)
	return ml
}

func (mp *MapPLZ) Add_LngLatPoly_Json(path [][]float64, props string) MapItem {
	var prop_map = map[string]interface{}{}
	json.Unmarshal([]byte(props), &prop_map)
	return mp.Add_LngLatPoly_Properties(path, prop_map)
}
