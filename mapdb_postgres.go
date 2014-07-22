package mapplz

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
)

type PSQLDatabase struct {
	db *sql.DB
}

func NewPostGISDB() *PSQLDatabase {
	return &PSQLDatabase{}
}

func (psql *PSQLDatabase) Type() string {
	return "postgis"
}

func (psql *PSQLDatabase) SetDB(dbinfo interface{}) {
	psql.db = dbinfo.(*sql.DB)
}

func (psql *PSQLDatabase) Add(mip MapItem) {
	var id int
	props_json, _ := json.Marshal(mip.Properties())
	props_str := string(props_json)
	wkt := mip.ToWKT()

	psql.db.QueryRow("INSERT INTO mapplz (properties, geom) VALUES ('" + props_str + "', ST_GeomFromText('" + wkt + "')) RETURNING id").Scan(&id)
}

func (psql *PSQLDatabase) Query() []MapItem {
	var mitems []MapItem
	rows, _ := psql.db.Query("SELECT id, ST_AsGeoJSON(geom) AS geo, properties FROM mapplz")

	defer rows.Close()
	for rows.Next() {
		var id int
		var geo string
		var props string
		if err := rows.Scan(&id, &geo, &props); err != nil {
			fmt.Printf("row scan error: %s", err)
		}
		faker := NewMapPLZ()
		mip := faker.Add(`{ "type": "Feature", "geometry": ` + geo + `}`)
		mitems = append(mitems, mip)
	}
	return mitems
}

func (psql *PSQLDatabase) Count() int {
  var count int
  rows, _ := psql.db.QueryRow("SELECT COUNT(*) FROM mapplz").Scan(&count)
  return count
}
