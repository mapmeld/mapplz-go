package mapplz

import (
	"encoding/json"
	"github.com/kellydunn/golang-geo"
	gj "github.com/kpawlik/geojson"
)

type MapItemPoly struct {
	path       *geo.Polygon
	properties map[string]interface{}
}

func NewMapItemPoly(latlngs [][]float64) *MapItemPoly {
	var polypts = []*geo.Point{}
	for i := 0; i < len(latlngs); i++ {
		polypts = append(polypts, geo.NewPoint(latlngs[i][0], latlngs[i][1]))
	}
	poly := geo.NewPolygon(polypts)
	return &MapItemPoly{path: poly, properties: make(map[string]interface{})}
}

func (mip *MapItemPoly) Type() string {
	return "polygon"
}

func (mip *MapItemPoly) Lat() float64 {
	return 0
}

func (mip *MapItemPoly) Lng() float64 {
	return 0
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

func (mip *MapItemPoly) SetProperties(props map[string]interface{}) {
	for key, value := range props {
		mip.properties[key] = value
	}
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
	return "POLYGON((0 0, 1 1))"
}
