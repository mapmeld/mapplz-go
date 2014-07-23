package mapplz

import (
	"database/sql"
	"fmt"
	"strings"
	_ "github.com/lib/pq"
)

type PSQLDatabase struct {
	db *sql.DB
}

func NewPostGISDB(db *sql.DB) *PSQLDatabase {
	return &PSQLDatabase{db: db}
}

func (psql *PSQLDatabase) Type() string {
	return "postgis"
}

func (psql *PSQLDatabase) QueryRow(sql string) int {
	var id int
	psql.db.QueryRow(sql).Scan(&id)
	return id
}

func (psql *PSQLDatabase) Query(sql string) []MapItem {
	if sql == "" {
		sql = "SELECT id, ST_AsGeoJSON(geom) AS geo, properties FROM mapplz"
	} else {
		sql_prop := strings.Split(strings.TrimSpace(sql), " ")[0]
		sql = strings.Replace(sql, sql_prop, "json_extract_path_text(properties, '" + sql_prop + "')", -1)
		sql = "SELECT id, ST_AsGeoJSON(geom) AS geo, properties FROM mapplz WHERE " + sql
	}

	var mitems []MapItem
	rows, _ := psql.db.Query(sql)

	defer rows.Close()
	for rows.Next() {
		var id int
		var geo string
		var props string
		if err := rows.Scan(&id, &geo, &props); err != nil {
			fmt.Printf("row scan error: %s", err)
		}
		mip := ConvertGeojsonFeature(`{ "type": "Feature", "geometry": `+geo+`}`, nil)
		mip.SetID(fmt.Sprintf("%v", id))
		mip.SetJsonProperties(props)
		mip.SetDB(psql)

		mitems = append(mitems, mip)
	}
	return mitems
}

func (psql *PSQLDatabase) Count(sql string) int {
	if sql == "" {
		sql = "SELECT COUNT(*) FROM mapplz"
	} else {
		sql_prop := strings.Split(strings.TrimSpace(sql), " ")[0]
		sql = strings.Replace(sql, sql_prop, "json_extract_path_text(properties, '" + sql_prop + "')", -1)
		sql = "SELECT COUNT(*) FROM mapplz WHERE " + sql
	}
	var count int
	psql.db.QueryRow(sql).Scan(&count)
	return count
}
