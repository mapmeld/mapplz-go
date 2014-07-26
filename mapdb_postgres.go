package mapplz

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"strings"
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

func (psql *PSQLDatabase) QueryRow(sql string) string {
	var id int
	psql.db.QueryRow(sql).Scan(&id)
	return fmt.Sprintf("%v", id)
}

func (psql *PSQLDatabase) Delete(id string) {
	sql := "DELETE from mapplz WHERE id = " + id
	psql.QueryRow(sql)
}

func (psql *PSQLDatabase) Save(sql interface{}) string {
	return psql.QueryRow(sql.(string))
}

func (psql *PSQLDatabase) Query(query interface{}) []MapItem {
	sql := query.(string)
	if sql == "" {
		sql = "SELECT id, ST_AsGeoJSON(geom) AS geo, properties FROM mapplz"
	} else {
		sql_prop := strings.Split(strings.TrimSpace(sql), " ")[0]
		sql = strings.Replace(sql, sql_prop, "json_extract_path_text(properties, '"+sql_prop+"')", -1)
		sql = "SELECT id, ST_AsGeoJSON(geom) AS geo, properties FROM mapplz WHERE " + sql
	}

	rows, _ := psql.db.Query(sql)
	return psql.responses(rows)
}

func (psql *PSQLDatabase) Count(query interface{}) int {
	sql := query.(string)
	if sql == "" {
		sql = "SELECT COUNT(*) FROM mapplz"
	} else {
		sql_prop := strings.Split(strings.TrimSpace(sql), " ")[0]
		sql = strings.Replace(sql, sql_prop, "json_extract_path_text(properties, '"+sql_prop+"')", -1)
		sql = "SELECT COUNT(*) FROM mapplz WHERE " + sql
	}
	var count int
	psql.db.QueryRow(sql).Scan(&count)
	return count
}

func (psql *PSQLDatabase) Within(area [][]float64) []MapItem {
	within_wkt := ""
	for i := 0; i < len(area); i++ {
		if i != 0 {
			within_wkt += ","
		}
		within_wkt += fmt.Sprintf("%v %v", area[i][1], area[i][0])
	}
	within_wkt = "POLYGON((" + within_wkt + "))"

	rows, _ := psql.db.Query("SELECT id, ST_AsGeoJSON(geom) AS geo, properties FROM mapplz AS start WHERE ST_Contains(ST_GeomFromText('" + within_wkt + "'), start.geom)")
	return psql.responses(rows)
}

func (psql *PSQLDatabase) Near(center []float64, count int) []MapItem {
	near_wkt := fmt.Sprintf("POINT(%v %v)", center[1], center[0])
	count_str := fmt.Sprintf("%v", count)

	var responses = []MapItem{}
	rows, _ := psql.db.Query("SELECT id, ST_AsGeoJSON(geom) AS geo, properties, ST_Distance(start.geom::geography, ST_GeomFromText('" + near_wkt + "')::geography) AS distance FROM mapplz AS start ORDER BY distance LIMIT " + count_str)

	defer rows.Close()
	for rows.Next() {
		var id int
		var geo string
		var props string
		var distance float64
		if err := rows.Scan(&id, &geo, &props, &distance); err != nil {
			fmt.Printf("row scan error: %s", err)
		}
		mip := ConvertGeojsonFeature(`{ "type": "Feature", "geometry": `+geo+`}`, nil)
		mip.SetID(fmt.Sprintf("%v", id))
		mip.SetJsonProperties(props)
		mip.SetDB(psql)

		responses = append(responses, mip)
	}
	return responses
}

func (psql *PSQLDatabase) responses(rows *sql.Rows) []MapItem {
	responses := []MapItem{}
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

		responses = append(responses, mip)
	}
	return responses
}
