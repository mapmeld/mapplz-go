package mapplz

import (
	"fmt"
	"testing"
)

func TestMapstore(t *testing.T) {
	mapstore := NewMapPLZ()
	mapstore.Add_Lat_Lng(40, -70)
	if mapstore.MapItems[0].Lat() != 40 {
		t.Errorf("point not made")
	} else {
		fmt.Printf("mapstore working")
	}
}

func TestGlobalAdd(t *testing.T) {
	mapstore := NewMapPLZ()

	// single param
	mapstore.Add([]float64{40, -70})
	mapstore.Add(`{ "type": "Feature", "geometry": { "type": "Point", "coordinates": [-70, 40] } }`)

	if len(mapstore.MapItems) != 2 || mapstore.Count("") != 2 {
		t.Errorf("global add failed")
	}

	// two param
	mapstore.Add2(40.121, -70.2)
	props := make(map[string]interface{})
	props["color"] = "red"
	mapstore.Add2([]float64{40.121, -70.2}, props)
	mapstore.Add2([]float64{40, -70}, `{ "color": "red" }`)

	if len(mapstore.MapItems) != 5 || mapstore.Count("") != 5 {
		t.Errorf("global add2 failed")
	}

	// three param
	props = make(map[string]interface{})
	props["color"] = "red"
	mapstore.Add3(40.121, -70.2, props)
	mapstore.Add3(40, -70, `{ "color": "red" }`)

	if len(mapstore.MapItems) != 7 || mapstore.Count("") != 7 {
		t.Errorf("global add3 failed")
	}
}

func TestLatLngProperties(t *testing.T) {
	mapstore := NewMapPLZ()
	props := make(map[string]interface{})
	props["color"] = "red"
	pt := mapstore.Add_Lat_Lng_Properties(40, -70, props)
	if pt.Properties()["color"] != "red" {
		t.Errorf("pt property not set")
	}
}

func TestLngLatJson(t *testing.T) {
	mapstore := NewMapPLZ()
	pt := mapstore.Add_Lng_Lat_Json(-70, 40, `{ "color": "#f00" }`)
	if pt.Properties()["color"] != "#f00" {
		t.Errorf("pt property not set")
	}
}

func TestLngLatParams(t *testing.T) {
	mapstore := NewMapPLZ()
	pt := mapstore.Add_Lng_Lat(-70, 40)
	if pt.Lat() != 40 || pt.Lng() != -70 {
		t.Errorf("point not made")
	}
}

func TestLatLngArray(t *testing.T) {
	mapstore := NewMapPLZ()
	latlng := []float64{40, -70}
	pt := mapstore.Add_LatLng(latlng)
	if pt.Lat() != 40 || pt.Lng() != -70 {
		t.Errorf("point not made")
	}
}

func TestLngLatArray(t *testing.T) {
	mapstore := NewMapPLZ()
	latlng := []float64{-70, 40}
	pt := mapstore.Add_LngLat(latlng)
	if pt.Lat() != 40 || pt.Lng() != -70 {
		t.Errorf("point not made")
	}
}

func TestLatLngPath(t *testing.T) {
	mapstore := NewMapPLZ()
	linepts := [][]float64{{40, -70}, {23.2, -110}}
	line := mapstore.Add_LatLngPath(linepts)
	first_pt := line.Path()[0][0]
	if first_pt[0] != 40 || first_pt[1] != -70 {
		t.Errorf("line not made from latlng path")
	}
}

func TestLngLatPath(t *testing.T) {
	mapstore := NewMapPLZ()
	linepts := [][]float64{{-70, 40}, {-110, 23.2}}
	line := mapstore.Add_LngLatPath(linepts)
	first_pt := line.Path()[0][0]
	if first_pt[0] != 40 || first_pt[1] != -70 {
		t.Errorf("line not made from lnglat path")
	}
}

func TestLngLatPathJson(t *testing.T) {
	mapstore := NewMapPLZ()
	linepts := [][]float64{{-70, 40}, {-110, 23.2}}

	line := mapstore.Add_LngLatPath_Json(linepts, `{ "color": "#f00" }`)
	if line.Properties()["color"] != "#f00" {
		t.Errorf("properties not added to lnglat path")
	}
}

func TestLatLngPoly(t *testing.T) {
	mapstore := NewMapPLZ()
	linepts := [][]float64{{40, -70}, {23.2, -110}, {25.2, -110}, {42.2, -70}, {40, -70}}
	line := mapstore.Add_LatLngPoly(linepts)
	first_pt := line.Path()[0][0]
	if first_pt[0] != 40 || first_pt[1] != -70 {
		t.Errorf("line not made from latlng path")
	}
}

func TestLngLatPoly(t *testing.T) {
	mapstore := NewMapPLZ()
	linepts := [][]float64{{-70, 40}, {-110, 23.2}, {-110, 25.2}, {-70, 42.2}, {-70, 40}}
	line := mapstore.Add_LngLatPoly(linepts)
	first_pt := line.Path()[0][0]
	if first_pt[0] != 40 || first_pt[1] != -70 {
		t.Errorf("line not made from lnglat path")
	}
}

func TestGeojsonPoint(t *testing.T) {
	mapstore := NewMapPLZ()
	gj := `{ "type": "Feature", "geometry": { "type": "Point", "coordinates": [-70, 40] } }`
	pt := mapstore.Add_Geojson_Feature(gj)

	if pt.Lat() != 40 || pt.Lng() != -70 {
		t.Errorf("geojson point not made")
	}
}

func TestGeojsonProperties(t *testing.T) {
	mapstore := NewMapPLZ()
	gj := `{ "type": "Feature", "geometry": { "type": "Point", "coordinates": [-70, 40] }, "properties": { "color": "#0f0" } }`
	pt := mapstore.Add_Geojson_Feature(gj)

	if pt.Properties()["color"] != "#0f0" {
		t.Errorf("geojson property not saved")
	}
}

func TestGeojsonLine(t *testing.T) {
	mapstore := NewMapPLZ()
	gj := `{ "type": "Feature", "geometry": { "type": "LineString", "coordinates": [[-70, 40], [-110, 32.1]] } }`
	line := mapstore.Add_Geojson_Feature(gj)

	first_pt := line.Path()[0][0]
	if first_pt[0] != 40 || first_pt[1] != -70 {
		t.Errorf("geojson line not made")
	}
}

func TestGeojsonPoly(t *testing.T) {
	mapstore := NewMapPLZ()
	gj := `{ "type": "Feature", "geometry": { "type": "Polygon", "coordinates": [[[-70, 40], [-110, 32.1], [-90, 25], [-70, 40]]] } }`
	poly := mapstore.Add_Geojson_Feature(gj)

	first_pt := poly.Path()[0][0]
	if first_pt[0] != 40 || first_pt[1] != -70 {
		t.Errorf("geojson poly not made")
	}
}

func TestGeojsonFeatureCollection(t *testing.T) {
	mapstore := NewMapPLZ()
	gj := `{ "type": "FeatureCollection", "features": [{ "type": "Feature", "geometry": { "type": "Point", "coordinates": [-70, 40] } }]}`
	gj_fc := mapstore.Add_Geojson_Collection(gj)

	if gj_fc[0].Lng() != -70 || gj_fc[0].Lat() != 40 {
		t.Errorf("geojson featurecollection not made")
	}
}

func TestGeojsonPointExport(t *testing.T) {
	mapstore := NewMapPLZ()
	pt := mapstore.Add_Lat_Lng(40, -70)
	output := pt.ToGeoJson()
	if output != `{"type":"Feature","geometry":{"type":"Point","coordinates":[-70,40]},"properties":null}` {
		t.Errorf("geojson output for point did not match")
	}
}

func TestGeojsonLineExport(t *testing.T) {
	mapstore := NewMapPLZ()
	linepts := [][]float64{{-70, 40}, {-110, 23.2}}
	line := mapstore.Add_LngLatPath(linepts)
	output := line.ToGeoJson()
	if output != `{"type":"Feature","geometry":{"type":"LineString","coordinates":[[-70,40],[-110,23.2]]},"properties":null}` {
		t.Errorf("geojson output for line did not match")
	}
}

func TestGeojsonPolyExport(t *testing.T) {
	mapstore := NewMapPLZ()
	polypts := [][]float64{{-70, 40}, {-110, 23.2}, {-97, 20}, {-70, 40}}
	line := mapstore.Add_LngLatPoly(polypts)
	output := line.ToGeoJson()
	if output != `{"type":"Feature","geometry":{"type":"Polygon","coordinates":[[[-70,40],[-110,23.2],[-97,20],[-70,40]]]},"properties":null}` {
		t.Errorf("geojson output for line did not match")
	}
}

func TestGeojsonAllExport(t *testing.T) {
	mapstore := NewMapPLZ()
	mapstore.Add_Lat_Lng(40, -70)
	output := mapstore.ToGeoJson()
	if output != `{"type":"FeatureCollection","features":[{"type":"Feature","geometry":{"type":"Point","coordinates":[-70,40]},"properties":null}]}` {
		t.Errorf("geojson output for all MapItems did not match")
	}
}

func TestWithin(t *testing.T) {
	mapstore := NewMapPLZ()
	mapstore.Add2(40.1, -70.2)
	mapstore.Add2(-40, -70)

	pts := mapstore.Within(`{ "type": "Feature", "geometry": { "type": "Polygon", "coordinates": [[[-71, 39], [-71, 41], [-69, 41], [-69, 39], [-71, 39]]] }}`)

	if len(pts) != 1 || pts[0].Lat() != 40.1 || pts[0].Lng() != -70.2 {
		t.Errorf("did not return point Within GeoJSON")
	}

	box := [][]float64{{39, -71}, {41, -71}, {41, -69}, {39, -69}, {39, -71}}
	pts = mapstore.Within(box)

	if len(pts) != 1 || pts[0].Lat() != 40.1 || pts[0].Lng() != -70.2 {
		t.Errorf("did not return point Within Array")
	}
}

func TestNear(t *testing.T) {
	mapstore := NewMapPLZ()
	mapstore.Add2(40.1, -70.2)
	mapstore.Add2(-40, -70)

	pts := mapstore.Near(`{ "type": "Feature", "geometry": { "type": "Point", "coordinates": [-71, 39] }}`, 1)

	if len(pts) != 1 || pts[0].Lat() != 40.1 || pts[0].Lng() != -70.2 {
		t.Errorf("did not return point from Near GeoJSON")
	}

	pt := []float64{-39, -71}
	pts = mapstore.Near(pt, 1)

	if len(pts) != 1 || pts[0].Lat() != -40 || pts[0].Lng() != -70 {
		t.Errorf("did not return point from Near Array")
	}
}

func TestDelete(t *testing.T) {
	mapstore := NewMapPLZ()
	pt := mapstore.Add_Lng_Lat(-70, 40)
	pt.Delete()
	if mapstore.Count("") != 0 {
		t.Errorf("did not delete point from nodb MapItems")
	}
}
