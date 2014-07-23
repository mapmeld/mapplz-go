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
