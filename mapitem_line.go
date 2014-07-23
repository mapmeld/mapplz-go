package mapplz

import (
	"encoding/json"
	"fmt"
	"github.com/kellydunn/golang-geo"
	gj "github.com/kpawlik/geojson"
)

type MapItemLine struct {
	id         string
	db         MapDatabase
	path       *geo.Polygon
	properties map[string]interface{}
}

func NewMapItemLine(latlngs [][]float64, db MapDatabase) *MapItemLine {
	var linepts = []*geo.Point{}
	for i := 0; i < len(latlngs); i++ {
		linepts = append(linepts, geo.NewPoint(latlngs[i][0], latlngs[i][1]))
	}
	line := geo.NewPolygon(linepts)
	return &MapItemLine{path: line, properties: make(map[string]interface{}), db: db}
}

func (mip *MapItemLine) Type() string {
	return "line"
}

func (mip *MapItemLine) Lat() float64 {
	return 0
}

func (mip *MapItemLine) Lng() float64 {
	return 0
}

func (mip *MapItemLine) Path() [][][]float64 {
	var path = [][][]float64{}
	var internic = [][]float64{}
	path = append(path, internic)
	path_pts := mip.path.Points()
	for i := 0; i < len(path_pts); i++ {
		var pt_entry = []float64{}
		pt_entry = append(pt_entry, path_pts[i].Lat())
		pt_entry = append(pt_entry, path_pts[i].Lng())
		path[0] = append(path[0], pt_entry)
	}
	return path
}

func (mip *MapItemLine) SetID(id string) {
	mip.id = id
}

func (mip *MapItemLine) SetDB(db MapDatabase) {
	mip.db = db
}

func (mip *MapItemLine) SetProperties(props map[string]interface{}) {
	for key, value := range props {
		mip.properties[key] = value
	}
	mip.Save()
}

func (mip *MapItemLine) SetJsonProperties(props string) {
	var prop_map = map[string]interface{}{}
	json.Unmarshal([]byte(props), &prop_map)
	mip.SetProperties(prop_map)
}

func (mip *MapItemLine) Properties() map[string]interface{} {
	return mip.properties
}

func (mip *MapItemLine) ToGeoJson() string {
	path_pts := mip.path.Points()
	var coords = []gj.Coordinate{}

	for i := 0; i < len(path_pts); i++ {
		lng := gj.CoordType(path_pts[i].Lng())
		lat := gj.CoordType(path_pts[i].Lat())
		coords = append(coords, gj.Coordinate{lng, lat})
	}

	gj_line := gj.NewLineString(gj.Coordinates(coords))
	feature := gj.NewFeature(gj_line, nil, nil)

	gjstr, err := gj.Marshal(feature)
	if err != nil {
		panic("failed to export point to GeoJSON")
	}
	return gjstr
}

func (mip *MapItemLine) ToWKT() string {
	path_pts := mip.path.Points()
	ptlist := ""
	for i := 0; i < len(path_pts); i++ {
		if i > 0 {
			ptlist += ","
		}
		ptlist += fmt.Sprintf("%v %v", path_pts[i].Lng(), path_pts[i].Lat())
	}
	return "LINESTRING(" + ptlist + ")"
}

func (mip *MapItemLine) Save() {
	if mip.db != nil {
		props_json, _ := json.Marshal(mip.Properties())
		props_str := string(props_json)
		wkt := mip.ToWKT()

		if mip.id == "" {
			// new MapItem
			id := mip.db.QueryRow("INSERT INTO mapplz (properties, geom) VALUES ('" + props_str + "', ST_GeomFromText('" + wkt + "')) RETURNING id")
			mip.id = fmt.Sprintf("%v", id)
		} else {
			// update MapItem
			mip.db.QueryRow("UPDATE mapplz SET geom = ST_GeomFromText('" + wkt + "'), properties = '" + props_str + "' WHERE id = " + mip.id)
		}
	}
}
