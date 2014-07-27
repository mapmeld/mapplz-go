package mapplz

import (
	"encoding/json"
	gj "github.com/mapmeld/geojson-bson"
	"sort"
)

type MapPLZ struct {
	MapItems []MapItem
	Database MapDatabase
}

// define MapItemList for custom distance sorts
type MapItemVector []MapItem
type MapItemList struct {
	Centroid []float64
	MIV      MapItemVector
}

func (v MapItemList) Swap(i, j int) {
	v.MIV[i], v.MIV[j] = v.MIV[j], v.MIV[i]
}
func (v MapItemList) Len() int {
	return len(v.MIV)
}
func (v MapItemList) Less(i, j int) bool {
	return v.MIV[i].DistanceFrom(v.Centroid) < v.MIV[j].DistanceFrom(v.Centroid)
}

func NewMapPLZ() MapPLZ {
	var mis = []MapItem{}
	return MapPLZ{MapItems: mis}
}

type MapDatabase interface {
	Type() string
	Query(interface{}) []MapItem
	Count(interface{}) int
	Within([][]float64) []MapItem
	Near([]float64, int) []MapItem
	QueryRow(string) string
	Save(interface{}) string
	Delete(string)
}

type MapItem interface {
	SetID(string)
	SetDB(MapDatabase)
	Type() string
	Lat() float64
	Lng() float64
	Path() [][][]float64
	SetProperties(map[string]interface{})
	Properties() map[string]interface{}
	SetJsonProperties(string)
	ToGeoJson() string
	ToGeoJsonFeature() *gj.Feature
	ToWKT() string
	Save()
	Delete()
	Deleted() bool
	Within([][]float64) bool
	DistanceFrom([]float64) float64
}

// global add

func (mp *MapPLZ) Add(input interface{}) MapItem {
	str, ok := input.(string)
	if ok {
		return mp.Add_Geojson_Feature(str)
	} else {
		arr, ok := input.([]interface{})
		if ok && len(arr) > 2 {
			var lat float64
			var lng float64
			var lat_int int
			var lng_int int

			lat_set, ok := arr[0].(float64)
			if !ok {
				lat_int, ok = arr[0].(int)
				lat = float64(lat_int)
			} else {
				lat = lat_set
			}

			lng_set, ok2 := arr[1].(float64)
			if !ok2 {
				lng_int, ok2 = arr[1].(int)
				lng = float64(lng_int)
			} else {
				lng = lng_set
			}

			props, ok := arr[2].(string)
			if ok {
				return mp.Add_Lat_Lng_Json(lat, lng, props)
			} else {
				props := arr[2].(map[string]interface{})
				return mp.Add_Lat_Lng_Properties(lat, lng, props)
			}
		} else {
			latlng, ok := input.([]float64)
			if ok {
  			return mp.Add_LatLng(latlng)
			} else {
				props := input.(map[string]interface{})
				if props["lat"] != nil && props["lng"] != nil {
					var lat float64
					var lng float64
					var lat_int int
					var lng_int int

					lat_set, ok := props["lat"].(float64)
					if !ok {
						lat_int, ok = props["lat"].(int)
						lat = float64(lat_int)
					} else {
						lat = lat_set
					}

					lng_set, ok2 := props["lng"].(float64)
					if !ok2 {
						lng_int, ok2 = props["lng"].(int)
						lng = float64(lng_int)
					} else {
						lng = lng_set
					}

					return mp.Add_Lat_Lng_Properties(lat, lng, props)
				} else {
					if props["path"] != nil {
						return mp.Add_LatLngPath_Properties(props["path"].([][]float64), props)
					} else {
						return nil
					}
				}
			}
		}
	}
}

func (mp *MapPLZ) Add2(input_first interface{}, input_second interface{}) MapItem {
	// interface_int.(float64) fails
	// ints must be read as int first, and then converted to float64
	var lat float64
	var lng float64
	var lat_int int
	var lng_int int

	lat_set, ok := input_first.(float64)
	if !ok {
		lat_int, ok = input_first.(int)
		lat = float64(lat_int)
	} else {
		lat = lat_set
	}

	lng_set, ok2 := input_second.(float64)
	if !ok2 {
		lng_int, ok2 = input_second.(int)
		lng = float64(lng_int)
	} else {
		lng = lng_set
	}

	if ok && ok2 {
		return mp.Add_Lat_Lng(lat, lng)
	} else {
		latlng, ok := input_first.([]float64)
		props, ok := input_second.(string)
		if ok {
			return mp.Add_LatLng_Json(latlng, props)
		} else {
			props := input_second.(map[string]interface{})
			return mp.Add_LatLng_Properties(latlng, props)
		}
	}
}

func (mp *MapPLZ) Add3(input_first interface{}, input_second interface{}, input_third interface{}) MapItem {
	// interface_int.(float64) fails
	// ints must be read as int first, and then converted to float64
	var lat float64
	var lng float64
	var lat_int int
	var lng_int int

	lat_set, ok := input_first.(float64)
	if !ok {
		lat_int, ok = input_first.(int)
		lat = float64(lat_int)
	} else {
		lat = lat_set
	}

	lng_set, ok2 := input_second.(float64)
	if !ok2 {
		lng_int, ok2 = input_second.(int)
		lng = float64(lng_int)
	} else {
		lng = lng_set
	}

	if ok && ok2 {
		props, ok := input_third.(string)
		if ok {
			return mp.Add_Lat_Lng_Json(lat, lng, props)
		} else {
			props := input_third.(map[string]interface{})
			return mp.Add_Lat_Lng_Properties(lat, lng, props)
		}
	} else {
		return nil
	}
}

// lat, lng with variations

func (mp *MapPLZ) Add_Lat_Lng(lat float64, lng float64) MapItem {
	mip := Convert_Lat_Lng(lat, lng, mp.Database)
	mip.Save()
	mp.MapItems = append(mp.MapItems, mip)
	return mip
}

func (mp *MapPLZ) Add_Lat_Lng_Properties(lat float64, lng float64, props map[string]interface{}) MapItem {
	mip := Convert_Lat_Lng(lat, lng, mp.Database)
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
	ml := ConvertPath(path, mp.Database)
	ml.Save()
	mp.MapItems = append(mp.MapItems, ml)
	return ml
}

func (mp *MapPLZ) Add_LatLngPath_Properties(path [][]float64, props map[string]interface{}) MapItem {
	ml := ConvertPath(path, mp.Database)
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
	ml := ConvertPoly(path, mp.Database)
	ml.Save()
	mp.MapItems = append(mp.MapItems, ml)
	return ml
}

func (mp *MapPLZ) Add_LatLngPoly_Properties(path [][]float64, props map[string]interface{}) MapItem {
	ml := ConvertPoly(path, mp.Database)
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

// database queries

func (mp *MapPLZ) Count(query interface{}) int {
	if mp.Database != nil {
		return mp.Database.Count(query)
	} else {
		return len(mp.Query(query))
	}
}

func (mp *MapPLZ) Query(query interface{}) []MapItem {
	if mp.Database != nil {
		return mp.Database.Query(query)
	} else {
		notnilitems := mp.NotNil(mp.MapItems)
		if query != nil && query != "" {
			conditions := query.(map[string]interface{})
			for qk := range conditions {
				for i := 0; i < len(notnilitems); i++ {
					if notnilitems[i].Properties()[qk] != conditions[qk] {
  					notnilitems[i] = nil
					}
				}
				notnilitems = mp.NotNil(notnilitems)
			}
		}
		return notnilitems
	}
}

func (mp *MapPLZ) Where(query interface{}) []MapItem {
	return mp.Query(query)
}

func (mp *MapPLZ) NotNil(items []MapItem) []MapItem {
	var response []MapItem
	for i := 0; i < len(items); i++ {
		if items[i] != nil && !items[i].Deleted() {
			response = append(response, items[i])
		}
	}
	return response
}

func (mp *MapPLZ) Within(area interface{}) []MapItem {
	area_poly := [][]float64{}
	area_json, ok := area.(string)
	if ok {
		// area_json is a GeoJSON of the poly
		area_poly = ConvertGeojsonFeature(area_json, nil).Path()[0]
	} else {
		// area is a Path
		area_poly = area.([][]float64)
	}

	if mp.Database != nil {
		return mp.Database.Within(area_poly)
	} else {
		var responses = []MapItem{}
		for i := 0; i < len(mp.MapItems); i++ {
			if !mp.MapItems[i].Deleted() && mp.MapItems[i].Within(area_poly) {
				responses = append(responses, mp.MapItems[i])
			}
		}
		return responses
	}
}

func (mp *MapPLZ) Near(center interface{}, count int) []MapItem {
	area_pt := []float64{}
	area_json, ok := center.(string)
	if ok {
		// area_json is a GeoJSON of the point
		mv := ConvertGeojsonFeature(area_json, nil)
		area_pt = append(area_pt, mv.Lat())
		area_pt = append(area_pt, mv.Lng())
	} else {
		// area_pt is a coordinate
		area_pt = center.([]float64)
	}

	if mp.Database != nil {
		return mp.Database.Near(area_pt, count)
	} else {
		bydistancevector := MapItemVector(mp.NotNil(mp.MapItems))
		bydistance := MapItemList{MIV: bydistancevector, Centroid: area_pt}
		sort.Sort(bydistance)
		return bydistance.MIV[0:count]
	}
}
