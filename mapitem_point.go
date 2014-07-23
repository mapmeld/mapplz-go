package mapplz

import (
	"encoding/json"
	"fmt"
	"github.com/kellydunn/golang-geo"
	gj "github.com/kpawlik/geojson"
)

type MapItemPoint struct {
	id         int
	db         *MapDatabase
	point      *geo.Point
	properties map[string]interface{}
}

func NewMapItemPoint(lat float64, lng float64) *MapItemPoint {
	pt := geo.NewPoint(lat, lng)
	return &MapItemPoint{point: pt, properties: make(map[string]interface{})}
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
	}
}
