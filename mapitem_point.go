package mapplz

import (
	"encoding/json"
	"fmt"
	"github.com/kellydunn/golang-geo"
	gj "github.com/kpawlik/geojson"
)

type MapItemPoint struct {
	id         string
	db         MapDatabase
	point      *geo.Point
	properties map[string]interface{}
}

func NewMapItemPoint(lat float64, lng float64, db MapDatabase) *MapItemPoint {
	pt := geo.NewPoint(lat, lng)
	return &MapItemPoint{point: pt, properties: make(map[string]interface{}), db: db}
}

func (mip *MapItemPoint) Type() string {
	return "point"
}

func (mip *MapItemPoint) Lat() float64 {
	return mip.point.Lat()
}

func (mip *MapItemPoint) Lng() float64 {
	return mip.point.Lng()
}

func (mip *MapItemPoint) Path() [][][]float64 {
	var blank = [][][]float64{{{}}}
	return blank
}

func (mip *MapItemPoint) SetID(id string) {
	mip.id = id
}

func (mip *MapItemPoint) SetDB(db MapDatabase) {
	mip.db = db
}

func (mip *MapItemPoint) SetProperties(props map[string]interface{}) {
	for key, value := range props {
		mip.properties[key] = value
	}
	mip.Save()
}

func (mip *MapItemPoint) SetJsonProperties(props string) {
	var prop_map = map[string]interface{}{}
	json.Unmarshal([]byte(props), &prop_map)
	mip.SetProperties(prop_map)
}

func (mip *MapItemPoint) Properties() map[string]interface{} {
	return mip.properties
}

func (mip *MapItemPoint) ToGeoJson() string {
	lng := gj.CoordType(mip.Lng())
	lat := gj.CoordType(mip.Lat())
	pt := gj.NewPoint(gj.Coordinate{lng, lat})
	feature := gj.NewFeature(pt, nil, nil)
	gjstr, err := gj.Marshal(feature)
	if err != nil {
		panic("failed to export point to GeoJSON")
	}
	return gjstr
}

func (mip *MapItemPoint) ToWKT() string {
	return fmt.Sprintf("POINT(%v %v)", mip.Lng(), mip.Lat())
}

func (mip *MapItemPoint) Save() {
	if mip.db != nil {
		props_json, _ := json.Marshal(mip.Properties())
		props_str := string(props_json)
		wkt := mip.ToWKT()

		if mip.id == "" {
			// new MapItem
			id := mip.db.QueryRow("INSERT INTO mapplz (properties, geom) VALUES ('" + props_str + "', ST_GeomFromText('" + wkt + "')) RETURNING id")
			mip.id = string(id)

		} else {
			// update MapItem
			mip.db.QueryRow("UPDATE mapplz SET geom = ST_GeomFromText('" + wkt + "'), properties = '" + props_str + "' WHERE id = " + mip.id)
		}
	}
}
