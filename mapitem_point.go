package mapplz

import (
  "encoding/json"
  "github.com/kellydunn/golang-geo"
  gj "github.com/kpawlik/geojson"
)

type MapItemPoint struct {
  point        *geo.Point
  properties   json.RawMessage
}

func NewMapItemPoint(lat float64, lng float64) *MapItemPoint {
  pt := geo.NewPoint(lat, lng)
	return &MapItemPoint{point: pt}
}

func (mip MapItemPoint) Type() string {
  return "point"
}

func (mip MapItemPoint) Lat() float64 {
  return mip.point.Lat()
}

func (mip MapItemPoint) Lng() float64 {
  return mip.point.Lng()
}

func (mip MapItemPoint) Path() [][][]float64 {
  var blank = [][][]float64{{{}}}
  return blank
}

func (mip MapItemPoint) Properties() json.RawMessage {
  return mip.properties
}

func (mip MapItemPoint) ToGeoJson() string {
  lng := gj.CoordType(mip.Lng())
  lat := gj.CoordType(mip.Lat())
  pt := gj.NewPoint(gj.Coordinate{lng, lat})
  feature := gj.NewFeature(pt, nil, nil)
  gjstr, err := gj.Marshal(feature)
  if err != nil {
    panic("failed to export point to GeoJSON")
  }
  return gjstr;
}
