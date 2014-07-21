package mapplz

import (
  "testing"
  _ "github.com/lib/pq"
  "database/sql"
)

func TestSetDb(t *testing.T) {
  mapstore := NewMapPLZ()
  db, err := sql.Open("postgres", "user=postgres dbname=travis_postgis sslmode=verify-full")
  if err != nil {
    t.Errorf("did not connect to PostGIS")
  }
  mapstore.Database = NewPostGISDB()
  mapstore.Database.SetDB(db)
}


func TestSaveToDb(t *testing.T) {
  mapstore := NewMapPLZ()
  db, err := sql.Open("postgres", "user=postgres dbname=travis_postgis sslmode=verify-full")
  if err != nil {
    t.Errorf("did not connect to PostGIS")
  }
  db.Exec("CREATE TABLE mapplz (id SERIAL PRIMARY KEY, properties JSON, geom public.geometry)")
  mapstore.Database = NewPostGISDB()
  mapstore.Database.SetDB(db)

  mapstore.Add2(40.1, -70.2)
  pt := mapstore.Query()[0]
  if (pt.Lat() != 40.1 || pt.Lng() != -70.2) {
    t.Errorf("did not return point from PostGIS")
  }

  db.Exec("DROP TABLE mapplz")
}
