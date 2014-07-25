package mapplz

import (
	"encoding/json"
	"fmt"
	"github.com/kellydunn/golang-geo"
	gj "github.com/kpawlik/geojson"
	"math"
)

type MapItemPoly struct {
	id         string
	db         MapDatabase
	path       *geo.Polygon
	properties map[string]interface{}
}

func NewMapItemPoly(latlngs [][]float64, db MapDatabase) *MapItemPoly {
	var polypts = []*geo.Point{}
	for i := 0; i < len(latlngs); i++ {
		polypts = append(polypts, geo.NewPoint(latlngs[i][0], latlngs[i][1]))
	}
	poly := geo.NewPolygon(polypts)
	return &MapItemPoly{path: poly, properties: make(map[string]interface{}), db: db}
}

func (mip *MapItemPoly) Type() string {
	return "polygon"
}

func (mip *MapItemPoly) Lat() float64 {
	var latsum float64
	path := mip.Path()[0]
	latsum = 0
	for i := 0; i < len(path); i++ {
		latsum += path[i][0]
	}
	return latsum / float64(len(path))
}

func (mip *MapItemPoly) Lng() float64 {
	var lngsum float64
	path := mip.Path()[0]
	lngsum = 0
	for i := 0; i < len(path); i++ {
		lngsum += path[i][1]
	}
	return lngsum / float64(len(path))
}

func (mip *MapItemPoly) Path() [][][]float64 {
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

func (mip *MapItemPoly) SetID(id string) {
	mip.id = id
}

func (mip *MapItemPoly) SetDB(db MapDatabase) {
	mip.db = db
}

func (mip *MapItemPoly) SetProperties(props map[string]interface{}) {
	for key, value := range props {
		mip.properties[key] = value
	}
	mip.Save()
}

func (mip *MapItemPoly) SetJsonProperties(props string) {
	var prop_map = map[string]interface{}{}
	json.Unmarshal([]byte(props), &prop_map)
	mip.SetProperties(prop_map)
}

func (mip *MapItemPoly) Properties() map[string]interface{} {
	return mip.properties
}

func (mip *MapItemPoly) ToGeoJson() string {
	path_pts := mip.path.Points()
	var coords = []gj.Coordinate{}

	for i := 0; i < len(path_pts); i++ {
		lng := gj.CoordType(path_pts[i].Lng())
		lat := gj.CoordType(path_pts[i].Lat())
		coords = append(coords, gj.Coordinate{lng, lat})
	}

	gj_poly := gj.NewPolygon(gj.MultiLine{gj.Coordinates(coords)})
	feature := gj.NewFeature(gj_poly, nil, nil)

	gjstr, err := gj.Marshal(feature)
	if err != nil {
		panic("failed to export point to GeoJSON")
	}
	return gjstr
}

func (mip *MapItemPoly) ToWKT() string {
	path_pts := mip.path.Points()
	ptlist := ""
	for i := 0; i < len(path_pts); i++ {
		if i > 0 {
			ptlist += ","
		}
		ptlist += fmt.Sprintf("%v %v", path_pts[i].Lng(), path_pts[i].Lat())
	}
	return "POLYGON((" + ptlist + "))"
}

func (mip *MapItemPoly) Save() {
	if mip.db != nil {
		props_json, _ := json.Marshal(mip.Properties())
		props_str := string(props_json)
		wkt := mip.ToWKT()

		if mip.id == "" {
			// new MapItem
			id := mip.db.Save("INSERT INTO mapplz (properties, geom) VALUES ('" + props_str + "', ST_GeomFromText('" + wkt + "')) RETURNING id")
			mip.id = id
		} else {
			// update MapItem
			mip.db.Save("UPDATE mapplz SET geom = ST_GeomFromText('" + wkt + "'), properties = '" + props_str + "' WHERE id = " + mip.id)
		}
	}
}

func (mip *MapItemPoly) DistanceFrom(center []float64) float64 {
	return geo.NewPoint(mip.Lat(), mip.Lng()).GreatCircleDistance(geo.NewPoint(center[0], center[1]))
}

func (mip *MapItemPoly) Within(area [][]float64) bool {
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
