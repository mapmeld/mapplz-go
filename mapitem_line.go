package mapplz

import (
	"github.com/kellydunn/golang-geo"
	gj "github.com/kpawlik/geojson"
)

type MapItemLine struct {
	path       *geo.Polygon
	properties map[string]interface{}
}

func NewMapItemLine(latlngs [][]float64) *MapItemLine {
	var linepts = []*geo.Point{}
	for i := 0; i < len(latlngs); i++ {
		linepts = append(linepts, geo.NewPoint(latlngs[i][0], latlngs[i][1]))
	}
	line := geo.NewPolygon(linepts)
	return &MapItemLine{path: line, properties: make(map[string]interface{})}
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

func (mip *MapItemLine) AddProperties(props map[string]interface{}) {
	for key, value := range props {
		mip.properties[key] = value
	}
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
