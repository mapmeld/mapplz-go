package mapplz

import (
	"gopkg.in/mgo.v2"
	"testing"
)

func TestSaveToMongo(t *testing.T) {
	mapstore := NewMapPLZ()
	session, err := mgo.Dial("localhost")
	defer session.Close()

	if err != nil {
		t.Errorf("did not connect to MongoDB")
	}

	collection := session.DB("sample").C("mapplz")
	collection.DropCollection()
	mapstore.Database = NewMongoDatabase(collection)

	mapstore.Add2(40.1, -70.2)
	pt := mapstore.Query(nil)[0]
	if pt.Lat() != 40.1 || pt.Lng() != -70.2 {
		t.Errorf("did not return point from MongoDB")
	}

	collection.DropCollection()
}

func TestMongoCount(t *testing.T) {
	mapstore := NewMapPLZ()
	session, err := mgo.Dial("localhost")
	defer session.Close()
	if err != nil {
		t.Errorf("did not connect to MongoDB")
	}
	collection := session.DB("sample").C("mapplz")
	mapstore.Database = NewMongoDatabase(collection)

	mapstore.Add3(40.1, -70.2, `{ "color": "red" }`)
	mapstore.Add3(40.1, -70.2, `{ "color": "blue" }`)

	if mapstore.Count(nil) != 2 {
		t.Errorf("did not count MongoDB points")
	}

	mquery := make(map[string]interface{})
	mquery["color"] = "blue"
	if mapstore.Count(mquery) != 1 {
		t.Errorf("did not filter MongoDB points in Count")
	}

	collection.DropCollection()
}

func TestMongoWhere(t *testing.T) {
	mapstore := NewMapPLZ()
	session, err := mgo.Dial("localhost")
	defer session.Close()
	if err != nil {
		t.Errorf("did not connect to MongoDB")
	}
	collection := session.DB("sample").C("mapplz")
	mapstore.Database = NewMongoDatabase(collection)

	mapstore.Add3(40.1, -70.2, `{ "color": "red" }`)
	mapstore.Add3(40.2, -70.3, `{ "color": "blue" }`)

	mquery := make(map[string]interface{})
	mquery["color"] = "blue"
	pt := mapstore.Where(mquery)[0]
	if pt.Lat() != 40.2 || pt.Lng() != -70.3 {
		t.Errorf("did not filter MongoDB points in Where")
	}

	collection.DropCollection()
}

func TestMongoUpdate(t *testing.T) {
	mapstore := NewMapPLZ()
	session, err := mgo.Dial("localhost")
	defer session.Close()
	if err != nil {
		t.Errorf("did not connect to MongoDB")
	}
	collection := session.DB("sample").C("mapplz")
	mapstore.Database = NewMongoDatabase(collection)

	pt := mapstore.Add2(40.1, -70.2)

	pt = mapstore.Query(nil)[0]
	if pt.Lat() != 40.1 || pt.Lng() != -70.2 {
		t.Errorf("did not return point from MongoDB")
	}

	props := make(map[string]interface{})
	props["color"] = "red"
	pt.SetProperties(props)

	if mapstore.Count(nil) != 1 {
		t.Errorf("did not keep to a single MongoDB point")
	}

	pt = mapstore.Query(nil)[0]
	if pt.Properties()["color"] != "red" {
		t.Errorf("did not update property in MongoDB")
	}

	collection.DropCollection()
}

func TestWithinMongo(t *testing.T) {
	mapstore := NewMapPLZ()
	session, err := mgo.Dial("localhost")
	defer session.Close()
	if err != nil {
		t.Errorf("did not connect to MongoDB")
	}
	collection := session.DB("sample").C("mapplz")
	geoindex := mgo.Index{Key: []string{"$2dsphere:geo.geometry"}, Bits: 26}
	collection.EnsureIndex(geoindex)
	mapstore.Database = NewMongoDatabase(collection)

	mapstore.Add2(40.1, -70.2)
	mapstore.Add2(-40, -70)

	pts := mapstore.Within(`{ "type": "Feature", "geometry": { "type": "Polygon", "coordinates": [[[-71, 39], [-71, 41], [-69, 41], [-69, 39], [-71, 39]]] }}`)

	if len(pts) != 1 || pts[0].Lat() != 40.1 || pts[0].Lng() != -70.2 {
		t.Errorf("did not return point from MongoDB Within GeoJSON")
	}

	box := [][]float64{{39, -71}, {41, -71}, {41, -69}, {39, -69}, {39, -71}}
	pts = mapstore.Within(box)

	if len(pts) != 1 || pts[0].Lat() != 40.1 || pts[0].Lng() != -70.2 {
		t.Errorf("did not return point from MongoDB Within Array")
	}

	collection.DropCollection()
}

func TestNearMongo(t *testing.T) {
	mapstore := NewMapPLZ()
	session, err := mgo.Dial("localhost")
	defer session.Close()
	if err != nil {
		t.Errorf("did not connect to MongoDB")
	}
	collection := session.DB("sample").C("mapplz")
	geoindex := mgo.Index{Key: []string{"$2dsphere:geo.geometry"}, Bits: 26}
	collection.EnsureIndex(geoindex)
	mapstore.Database = NewMongoDatabase(collection)

	mapstore.Add2(40.1, -70.2)
	mapstore.Add2(-40, -70)

	pts := mapstore.Near(`{ "type": "Feature", "geometry": { "type": "Point", "coordinates": [-71, 39] }}`, 1)

	if len(pts) != 1 || pts[0].Lat() != 40.1 || pts[0].Lng() != -70.2 {
		t.Errorf("did not return point from MongoDB Near GeoJSON")
	}

	pt := []float64{-39, -71}
	pts = mapstore.Near(pt, 1)

	if len(pts) != 1 || pts[0].Lat() != -40 || pts[0].Lng() != -70 {
		t.Errorf("did not return point from MongoDB Near Array")
	}

	collection.DropCollection()
}

func TestLngLatPathJsonMongo(t *testing.T) {
	mapstore := NewMapPLZ()
	session, err := mgo.Dial("localhost")
	defer session.Close()
	if err != nil {
		t.Errorf("did not connect to MongoDB")
	}
	collection := session.DB("sample").C("mapplz")
	geoindex := mgo.Index{Key: []string{"$2dsphere:geo.geometry"}, Bits: 26}
	collection.EnsureIndex(geoindex)
	mapstore.Database = NewMongoDatabase(collection)

	linepts := [][]float64{{-70, 40}, {-110, 23.2}}
	line := mapstore.Add_LngLatPath_Json(linepts, `{ "color": "#f00" }`)
	if line.Properties()["color"] != "#f00" {
		t.Errorf("properties not added to lnglat path on MongoDB")
	}

	collection.DropCollection()
}

func TestLatLngPolyMongo(t *testing.T) {
	mapstore := NewMapPLZ()
	session, err := mgo.Dial("localhost")
	defer session.Close()
	if err != nil {
		t.Errorf("did not connect to MongoDB")
	}
	collection := session.DB("sample").C("mapplz")
	geoindex := mgo.Index{Key: []string{"$2dsphere:geo.geometry"}, Bits: 26}
	collection.EnsureIndex(geoindex)
	mapstore.Database = NewMongoDatabase(collection)

	linepts := [][]float64{{40, -70}, {23.2, -110}, {25.2, -110}, {42.2, -70}, {40, -70}}
	line := mapstore.Add_LatLngPoly(linepts)
	first_pt := line.Path()[0][0]
	if first_pt[0] != 40 || first_pt[1] != -70 {
		t.Errorf("line not made from latlng path on MongoDB")
	}

	collection.DropCollection()
}
