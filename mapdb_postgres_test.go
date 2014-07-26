package mapplz

import (
	"database/sql"
	_ "github.com/lib/pq"
	"testing"
)

func TestSaveToDb(t *testing.T) {
	mapstore := NewMapPLZ()
	db, err := sql.Open("postgres", "user=postgres dbname=travis_postgis sslmode=disable")
	if err != nil {
		t.Errorf("did not connect to PostGIS")
	}
	db.Exec("CREATE TABLE mapplz (id SERIAL PRIMARY KEY, properties JSON, geom public.geometry)")

	mapstore.Database = NewPostGISDB(db)

	mapstore.Add2(40.1, -70.2)
	pt := mapstore.Query("")[0]
	if pt.Lat() != 40.1 || pt.Lng() != -70.2 {
		t.Errorf("did not return point from PostGIS")
	}

	db.Exec("DROP TABLE mapplz")
}

func TestCount(t *testing.T) {
	mapstore := NewMapPLZ()
	db, err := sql.Open("postgres", "user=postgres dbname=travis_postgis sslmode=disable")
	if err != nil {
		t.Errorf("did not connect to PostGIS")
	}
	db.Exec("CREATE TABLE mapplz (id SERIAL PRIMARY KEY, properties JSON, geom public.geometry)")

	mapstore.Database = NewPostGISDB(db)

	mapstore.Add3(40.1, -70.2, `{ "color": "red" }`)
	mapstore.Add3(40.1, -70.2, `{ "color": "blue" }`)
	if mapstore.Count("") != 2 {
		t.Errorf("did not count PostGIS points")
	}

	if mapstore.Count("color = 'blue'") != 1 {
		t.Errorf("did not filter PostGIS points in Count")
	}

	db.Exec("DROP TABLE mapplz")
}

func TestWhere(t *testing.T) {
	mapstore := NewMapPLZ()
	db, err := sql.Open("postgres", "user=postgres dbname=travis_postgis sslmode=disable")
	if err != nil {
		t.Errorf("did not connect to PostGIS")
	}
	db.Exec("CREATE TABLE mapplz (id SERIAL PRIMARY KEY, properties JSON, geom public.geometry)")

	mapstore.Database = NewPostGISDB(db)

	mapstore.Add3(40.1, -70.2, `{ "color": "red" }`)
	mapstore.Add3(40.2, -70.3, `{ "color": "blue" }`)

	pt := mapstore.Where("color = 'blue'")[0]
	if pt.Lat() != 40.2 || pt.Lng() != -70.3 {
		t.Errorf("did not filter PostGIS points in Where")
	}

	db.Exec("DROP TABLE mapplz")
}

func TestUpdate(t *testing.T) {
	mapstore := NewMapPLZ()
	db, err := sql.Open("postgres", "user=postgres dbname=travis_postgis sslmode=disable")
	if err != nil {
		t.Errorf("did not connect to PostGIS")
	}
	db.Exec("CREATE TABLE mapplz (id SERIAL PRIMARY KEY, properties JSON, geom public.geometry)")

	mapstore.Database = NewPostGISDB(db)

	pt := mapstore.Add2(40.1, -70.2)

	pt = mapstore.Query("")[0]
	if pt.Lat() != 40.1 || pt.Lng() != -70.2 {
		t.Errorf("did not return point from PostGIS")
	}

	props := make(map[string]interface{})
	props["color"] = "red"

	pt.SetProperties(props)

	pt = mapstore.Query("")[0]
	if pt.Properties()["color"] != "red" {
		t.Errorf("did not update property in PostGIS")
	}

	if mapstore.Count("") != 1 {
		t.Errorf("did not keep to a single PostGIS point")
	}

	db.Exec("DROP TABLE mapplz")
}

func TestWithinPSQL(t *testing.T) {
	mapstore := NewMapPLZ()
	db, err := sql.Open("postgres", "user=postgres dbname=travis_postgis sslmode=disable")
	if err != nil {
		t.Errorf("did not connect to PostGIS")
	}
	db.Exec("CREATE TABLE mapplz (id SERIAL PRIMARY KEY, properties JSON, geom public.geometry)")

	mapstore.Database = NewPostGISDB(db)

	mapstore.Add2(40.1, -70.2)
	mapstore.Add2(-40, -70)

	pts := mapstore.Within(`{ "type": "Feature", "geometry": { "type": "Polygon", "coordinates": [[[-71, 39], [-71, 41], [-69, 41], [-69, 39], [-71, 39]]] }}`)

	if len(pts) != 1 || pts[0].Lat() != 40.1 || pts[0].Lng() != -70.2 {
		t.Errorf("did not return point from PostGIS Within GeoJSON")
	}

	box := [][]float64{{39, -71}, {41, -71}, {41, -69}, {39, -69}, {39, -71}}
	pts = mapstore.Within(box)

	if len(pts) != 1 || pts[0].Lat() != 40.1 || pts[0].Lng() != -70.2 {
		t.Errorf("did not return point from PostGIS Within Array")
	}

	db.Exec("DROP TABLE mapplz")
}

func TestNearPSQL(t *testing.T) {
	mapstore := NewMapPLZ()
	db, err := sql.Open("postgres", "user=postgres dbname=travis_postgis sslmode=disable")
	if err != nil {
		t.Errorf("did not connect to PostGIS")
	}
	db.Exec("CREATE TABLE mapplz (id SERIAL PRIMARY KEY, properties JSON, geom public.geometry)")

	mapstore.Database = NewPostGISDB(db)

	mapstore.Add2(40.1, -70.2)
	mapstore.Add2(-40, -70)

	pts := mapstore.Near(`{ "type": "Feature", "geometry": { "type": "Point", "coordinates": [-71, 39] }}`, 1)

	if len(pts) != 1 || pts[0].Lat() != 40.1 || pts[0].Lng() != -70.2 {
		t.Errorf("did not return point from PostGIS Near GeoJSON")
	}

	pt := []float64{-39, -71}
	pts = mapstore.Near(pt, 1)

	if len(pts) != 1 || pts[0].Lat() != -40 || pts[0].Lng() != -70 {
		t.Errorf("did not return point from PostGIS Near Array")
	}

	db.Exec("DROP TABLE mapplz")
}

func TestLngLatPathJsonGIS(t *testing.T) {
	mapstore := NewMapPLZ()
	db, err := sql.Open("postgres", "user=postgres dbname=travis_postgis sslmode=disable")
	if err != nil {
		t.Errorf("did not connect to PostGIS")
	}
	db.Exec("CREATE TABLE mapplz (id SERIAL PRIMARY KEY, properties JSON, geom public.geometry)")
	mapstore.Database = NewPostGISDB(db)

	linepts := [][]float64{{-70, 40}, {-110, 23.2}}
	line := mapstore.Add_LngLatPath_Json(linepts, `{ "color": "#f00" }`)
	if line.Properties()["color"] != "#f00" {
		t.Errorf("properties not added to lnglat path on PostGIS")
	}

	db.Exec("DROP TABLE mapplz")
}

func TestLatLngPolyGIS(t *testing.T) {
	mapstore := NewMapPLZ()
	db, err := sql.Open("postgres", "user=postgres dbname=travis_postgis sslmode=disable")
	if err != nil {
		t.Errorf("did not connect to PostGIS")
	}
	db.Exec("CREATE TABLE mapplz (id SERIAL PRIMARY KEY, properties JSON, geom public.geometry)")
	mapstore.Database = NewPostGISDB(db)

	linepts := [][]float64{{40, -70}, {23.2, -110}, {25.2, -110}, {42.2, -70}, {40, -70}}
	line := mapstore.Add_LatLngPoly(linepts)
	first_pt := line.Path()[0][0]
	if first_pt[0] != 40 || first_pt[1] != -70 {
		t.Errorf("line not made from latlng path on PostGIS")
	}

	db.Exec("DROP TABLE mapplz")
}
