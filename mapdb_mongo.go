package mapplz

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"fmt"
)

type MongoDatabase struct {
	collection *mgo.Collection
}

func NewMongoDatabase(collection *mgo.Collection) *MongoDatabase {
	return &MongoDatabase{collection: collection}
}

func (mdb *MongoDatabase) Type() string {
	return "mongodb"
}

func (mdb *MongoDatabase) QueryRow(sql string) string {
  return "0"
}

func (mdb *MongoDatabase) Save(mquery interface{}) string {
	mdoc := mquery.(map[string]interface{})

	// return ID of new documents with fake Upsert
	// replace with a real ID for updates
	query_id := bson.M{"fake_field": "fake_value"}
	if mdoc["id"] != nil {
		query_id = bson.M{"_id": mdoc["id"]}
	}

	changes, err := mdb.collection.Upsert(query_id, bson.M(mdoc))
	if err != nil {
    panic(err)
  }

	new_id := changes.UpsertedId.(bson.ObjectId)

	return new_id.String()
}

func (mdb *MongoDatabase) Query(sql string) []MapItem {
	mitems := []MapItem{}

	var results []interface{}
	mdb.collection.Find(nil).All(&results)

	for i := 0; i < len(results); i++ {
		result_map := results[i].(bson.M)
		fmt.Printf("%s", result_map)
		mip := ConvertGeojsonFeature(result_map["geo"].(string), nil)
		bson_id := result_map["_id"].(bson.ObjectId)
		mip.SetID(bson_id.String())
		// mip.SetJsonProperties(result_map["props"].(string))
		mip.SetDB(mdb)
		mitems = append(mitems, mip)
	}
	return mitems
}

func (mdb *MongoDatabase) Count(sql string) int {
  return 0
}

func (mdb *MongoDatabase) Within(area [][]float64) []MapItem {
  return []MapItem{}
}

func (mdb *MongoDatabase) Near(center []float64, count int) []MapItem {
  return []MapItem{}
}
