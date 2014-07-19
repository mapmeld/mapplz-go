package mapplz

import (
	"testing"
)

func TestMapstore(t *testing.T) {
  mapstore := NewMapPLZ()
  mapstore.Add_Lat_Lng(40, -70)
  if(mapstore.MapItems[0].Lat() != 40) {
    t.Errorf("point not made")
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
	if(first_pt[0] != 40 || first_pt[1] != -70) {
		t.Errorf("line not made from latlng path")
	}
}

func TestLngLatPath(t *testing.T) {
	mapstore := NewMapPLZ()
	linepts := [][]float64{{-70, 40}, {-110, 23.2}}
	line := mapstore.Add_LngLatPath(linepts)
	first_pt := line.Path()[0][0]
	if(first_pt[0] != 40 || first_pt[1] != -70) {
		t.Errorf("line not made from lnglat path")
	}
}

func TestLatLngPoly(t *testing.T) {
	mapstore := NewMapPLZ()
	linepts := [][]float64{{40, -70}, {23.2, -110}}
	line := mapstore.Add_LatLngPoly(linepts)
	first_pt := line.Path()[0][0]
	if(first_pt[0] != 40 || first_pt[1] != -70) {
		t.Errorf("line not made from latlng path")
	}
}

func TestLngLatPoly(t *testing.T) {
	mapstore := NewMapPLZ()
	linepts := [][]float64{{-70, 40}, {-110, 23.2}}
	line := mapstore.Add_LngLatPoly(linepts)
	first_pt := line.Path()[0][0]
	if(first_pt[0] != 40 || first_pt[1] != -70) {
		t.Errorf("line not made from lnglat path")
	}
}

func TestGeojsonPoint(t *testing.T) {
	mapstore := NewMapPLZ()
	gj := []byte(`{ "type": "Feature", "geometry": { "type": "Point", "coordinates": [-70, 40] } }`)
	pt := mapstore.Add_Geojson_Feature(gj)

	if pt.Lat() != 40 || pt.Lng() != -70 {
		t.Errorf("geojson point not made")
	}
}

func TestGeojsonLine(t *testing.T) {
	mapstore := NewMapPLZ()
	gj := []byte(`{ "type": "Feature", "geometry": { "type": "LineString", "coordinates": [[-70, 40], [-110, 32.1]] } }`)
	line := mapstore.Add_Geojson_Feature(gj)

	first_pt := line.Path()[0][0]
	if(first_pt[0] != 40 || first_pt[1] != -70) {
		t.Errorf("geojson line not made")
	}
}

func TestGeojsonPoly(t *testing.T) {
	mapstore := NewMapPLZ()
  gj := []byte(`{ "type": "Feature", "geometry": { "type": "Polygon", "coordinates": [[[-70, 40], [-110, 32.1], [-90, 25], [-70, 40]]] } }`)
  poly := mapstore.Add_Geojson_Feature(gj)

  first_pt := poly.Path()[0][0]
	if(first_pt[0] != 40 || first_pt[1] != -70) {
		t.Errorf("geojson poly not made")
	}
}

func TestGeojsonFeatureCollection(t *testing.T) {
	mapstore := NewMapPLZ()
  gj := []byte(`{ "type": "FeatureCollection", "features": [{ "type": "Feature", "geometry": { "type": "Point", "coordinates": [-70, 40] } }]}`)
  gj_fc := mapstore.Add_Geojson_Collection(gj)

  if gj_fc.Features[0].Geometry.Point.Coordinates[0] != -70 || gj_fc.Features[0].Geometry.Point.Coordinates[1] != 40 {
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
