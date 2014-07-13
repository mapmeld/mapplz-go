package mapplz

import (
	"testing"
)

func TestLatLngParams(t *testing.T) {
	pt := Add_Lat_Lng(40, -70)
	if pt.Lat() != 40 || pt.Lng() != -70 {
		t.Errorf("point not made")
	}
}

func TestLngLatParams(t *testing.T) {
	pt := Add_Lng_Lat(-70, 40)
	if pt.Lat() != 40 || pt.Lng() != -70 {
		t.Errorf("point not made")
	}
}

func TestLatLngArray(t *testing.T) {
	latlng := []float64{40, -70}
	pt := Add_LatLng(latlng)
	if pt.Lat() != 40 || pt.Lng() != -70 {
		t.Errorf("point not made")
	}
}

func TestLngLatArray(t *testing.T) {
	latlng := []float64{-70, 40}
	pt := Add_LngLat(latlng)
	if pt.Lat() != 40 || pt.Lng() != -70 {
		t.Errorf("point not made")
	}
}

func TestGeojsonPoint(t *testing.T) {
	gj := []byte(`{ "type": "Feature", "geometry": { "type": "Point", "coordinates": [-70, 40] } }`)
	gj_pt := Add_Geojson_Feature_Str(gj)

	if gj_pt.Geometry.Point.Coordinates[0] != -70 || gj_pt.Geometry.Point.Coordinates[1] != 40 {
		t.Errorf("geojson point not made")
	}
}

func TestGeojsonLine(t *testing.T) {
	gj := []byte(`{ "type": "Feature", "geometry": { "type": "LineString", "coordinates": [[-70, 40], [-110, 32.1]] } }`)
	gj_line := Add_Geojson_Feature_Str(gj)

	if gj_line.Geometry.Line.Coordinates[0][0] != -70 || gj_line.Geometry.Line.Coordinates[0][1] != 40 {
		t.Errorf("geojson line not made")
	}
}
