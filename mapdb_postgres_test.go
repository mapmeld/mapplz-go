package mapplz

import (
	"database/sql"
	_ "github.com/lib/pq"
	"testing"
)

func TestSetDb(t *testing.T) {
	mapstore := NewMapPLZ()
	db, err := sql.Open("postgres", "user=postgres dbname=travis_postgis sslmode=disable")
	if err != nil {
		t.Errorf("did not connect to PostGIS")
	}
	mapstore.Database = NewPostGISDB()
	mapstore.Database.SetDB(db)
}

func TestSaveToDb(t *testing.T) {
	mapstore := NewMapPLZ()
	db, err := sql.Open("postgres", "user=postgres dbname=travis_postgis sslmode=disable")
	if err != nil {
		t.Errorf("did not connect to PostGIS")
	}
	db.Exec("CREATE TABLE mapplz (id SERIAL PRIMARY KEY, properties JSON, geom public.geometry)")

	mapstore.Database = NewPostGISDB()
	mapstore.Database.SetDB(db)

	mapstore.Add2(40.1, -70.2)
	pt := mapstore.Query()[0]
	if pt.Lat() != 40.1 || pt.Lng() != -70.2 {
		t.Errorf("did not return point from PostGIS")
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

	mapstore.Database = NewPostGISDB()
	mapstore.Database.SetDB(db)

	mapstore.Add2(40.1, -70.2)
	pt := mapstore.Query()[0]
	if pt.Lat() != 40.1 || pt.Lng() != -70.2 {
		t.Errorf("did not return point from PostGIS")
	}

	props := make(map[string]interface{})
	props["color"] = "red"
	pt.SetProperties(props)

	pt = mapstore.Query()[0]
	if pt.Properties()["color"] != "red" {
		t.Errorf("did not update property in PostGIS")
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

	mapstore.Database = NewPostGISDB()
	mapstore.Database.SetDB(db)

	mapstore.Add2(40.1, -70.2)
	if mapstore.Count() != 1 {
		t.Errorf("did not count PostGIS point")
	}

	db.Exec("DROP TABLE mapplz")
}
