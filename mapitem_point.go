package mapplz

import (
	"encoding/json"
	"fmt"
	"github.com/kellydunn/golang-geo"
	gj "github.com/mapmeld/geojson-bson"
	"math"
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

func (mip *MapItemPoint) ToGeoJsonFeature() *gj.Feature {
	lng := gj.CoordType(mip.Lng())
	lat := gj.CoordType(mip.Lat())
	pt := gj.NewPoint(gj.Coordinate{lng, lat})
	return gj.NewFeature(pt, nil, nil)
}

func (mip *MapItemPoint) ToGeoJson() string {
	feature := mip.ToGeoJsonFeature()
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
		if mip.id == "" {
			// new MapItem
			var id string
			if mip.db.Type() == "postgis" {
				props_json, _ := json.Marshal(mip.Properties())
				props_str := string(props_json)
			  id = mip.db.Save("INSERT INTO mapplz (properties, geom) VALUES ('" + props_str + "', ST_GeomFromText('" + mip.ToWKT() + "')) RETURNING id")
			} else {
				mdoc := make(map[string]interface{})
				props := mip.Properties()
				for key := range props {
					mdoc[key] = props[key]
				}
				mdoc["geo"] = mip.ToGeoJsonFeature()
				id = mip.db.Save(mdoc)
			}
			mip.id = id
		} else {
			// update MapItem
			if mip.db.Type() == "postgis" {
				props_json, _ := json.Marshal(mip.Properties())
				props_str := string(props_json)
  			mip.db.Save("UPDATE mapplz SET geom = ST_GeomFromText('" + mip.ToWKT() + "'), properties = '" + props_str + "' WHERE id = " + mip.id)
			} else {
				mdoc := make(map[string]interface{})
				mdoc["id"] = mip.id
				props := mip.Properties()
				for key := range props {
					mdoc[key] = props[key]
				}
				mdoc["geo"] = mip.ToGeoJsonFeature()
				mip.db.Save(mdoc)
			}
		}
	}
}

func (mip *MapItemPoint) DistanceFrom(center []float64) float64 {
	return mip.point.GreatCircleDistance(geo.NewPoint(center[0], center[1]))
}

func (mip *MapItemPoint) Within(area [][]float64) bool {
	// adapted from polyclip-go by Mateusz Czapliński
	intersections := 0
	for i := range area {
		curr := area[i]
		ii := i + 1
		if ii == len(area) {
			ii = 0
		}
		next := area[ii]

		if (mip.Lat() >= next[0] || mip.Lat() <= curr[0]) &&
			(mip.Lat() >= curr[0] || mip.Lat() <= next[0]) {
			continue
		}
		// Edge is from curr to next.

		if mip.Lng() >= math.Max(curr[1], next[1]) || next[0] == curr[0] {
			continue
		}

		// Find where the line intersects...
		xint := (mip.Lat()-curr[0])*(next[1]-curr[1])/(next[0]-curr[0]) + curr[1]
		if curr[1] != next[1] && mip.Lng() > xint {
			continue
		}

		intersections++
	}

	return (intersections%2 != 0)
}
