package mapplz

import (
	"testing"
	"gopkg.in/mgo.v2"
)

func TestSaveToMongo(t *testing.T) {
	mapstore := NewMapPLZ()
	session, err := mgo.Dial("localhost")
	defer session.Close()

	if err != nil {
		t.Errorf("did not connect to MongoDB")
	}

	collection := session.DB("sample").C("mapplz")

	mapstore.Database = NewMongoDatabase(collection)

	mapstore.Add2(40.1, -70.2)
	pt := mapstore.Query("")[0]
	if pt.Lat() != 40.1 || pt.Lng() != -70.2 {
		t.Errorf("did not return point from MongoDB")
	}

	collection.DropCollection()
}
