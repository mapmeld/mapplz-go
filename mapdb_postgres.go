package mapplz

import (
	"database/sql"
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

func (psql *PSQLDatabase) QueryRow(sql string) int {
	var id int
	psql.db.QueryRow(sql).Scan(&id)
	return id
}

func (psql *PSQLDatabase) SetDB(dbinfo interface{}) {
	psql.db = dbinfo.(*sql.DB)
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
		mip := ConvertGeojsonFeature(`{ "type": "Feature", "geometry": `+geo+`}`, nil)
		mip.SetJsonProperties(props)

		// set id and db later so it's not re-inserted into DB
		mip.SetID(string(id))
		mip.SetDB(psql)
		mitems = append(mitems, mip)
	}
	return mitems
}

func (psql *PSQLDatabase) Count() int {
	var count int
	psql.db.QueryRow("SELECT COUNT(*) FROM mapplz").Scan(&count)
	return count
}
