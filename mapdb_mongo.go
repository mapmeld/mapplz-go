package mapplz

import (
	"encoding/json"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"strings"
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

func (mdb *MongoDatabase) Delete(id string) {
	query_id := bson.M{"_id": bson.ObjectIdHex(id)}
	mdb.collection.Remove(query_id)
}

func (mdb *MongoDatabase) Save(mquery interface{}) string {
	mdoc := mquery.(map[string]interface{})

	// return ID of new documents with fake Upsert
	// replace with a real ID for updates
	var string_id string
	query_id := bson.M{"fake_field": "fake_value"}

	if mdoc["id"] != nil {
		query_id = bson.M{"_id": bson.ObjectIdHex(mdoc["id"].(string))}
		string_id = mdoc["id"].(string)
		mdoc["id"] = nil
		err := mdb.collection.Update(query_id, bson.M(mdoc))
		if err != nil {
			panic(err)
		}
	} else {
		changes, err := mdb.collection.Upsert(query_id, bson.M(mdoc))
		if err != nil {
			panic(err)
		}
		proposed_id, ok := changes.UpsertedId.(string)
		if ok {
			string_id = proposed_id
		} else {
			new_id := changes.UpsertedId.(bson.ObjectId)
			string_id = new_id.String()
		}
	}

	return mdb.sanitize(string_id)
}

func (mdb *MongoDatabase) Query(query interface{}) []MapItem {
	mdoc := bson.M{}
	if query != nil {
		// support empty string query
		_, ok := query.(string)
		if ok {
			query = nil
		} else {
			query_map := query.(map[string]interface{})
			for qk := range query_map {
				if qk != "id" && qk != "_id" {
					mdoc[qk] = query_map[qk]
				}
			}
		}
	}

	var results []interface{}
	mdb.collection.Find(mdoc).All(&results)

	mitems := []MapItem{}
	for i := 0; i < len(results); i++ {
		result_map := results[i].(bson.M)
		geo_str, _ := json.Marshal(result_map["geo"])
		mip := ConvertGeojsonFeature(string(geo_str), nil)

		string_id, ok := result_map["_id"].(string)
		if ok {
			mip.SetID(string_id)
		} else {
			bson_id := result_map["_id"].(bson.ObjectId)
			mip.SetID(mdb.sanitize(bson_id.String()))
		}

		props_map := make(map[string]interface{})
		for key := range result_map {
			if key != "_id" && key != "id" && key != "geo" {
				props_map[key] = result_map[key]
			}
		}
		mip.SetProperties(props_map)

		mip.SetDB(mdb)
		mitems = append(mitems, mip)
	}
	return mitems
}

func (mdb *MongoDatabase) Count(query interface{}) int {
	mquery := bson.M{}

	if query != nil {
		// support empty string query
		_, ok := query.(string)
		if ok {
			query = nil
		} else {
			mdoc := query.(map[string]interface{})
			for key := range mdoc {
				mquery[key] = mdoc[key]
			}
		}
	}

	count, _ := mdb.collection.Find(mquery).Count()
	return count
}

func (mdb *MongoDatabase) Within(area [][]float64) []MapItem {
	gjson := make(map[string]interface{})
	gjson["type"] = "Polygon"
	rev_path := [][][]float64{{{}}}
	for i := 0; i < len(area); i++ {
		lat := area[i][0]
		lng := area[i][1]
		area[i][0] = lng
		area[i][1] = lat
	}
	rev_path[0] = area
	gjson["coordinates"] = rev_path

	geometry := make(map[string]interface{})
	geometry["$geometry"] = gjson
	gw := make(map[string]interface{})
	gw["$geoWithin"] = geometry

	qi := make(map[string]interface{})
	qi["geo.geometry"] = gw

	return mdb.Query(qi)
}

func (mdb *MongoDatabase) Near(center []float64, count int) []MapItem {
	gjson := make(map[string]interface{})
	gjson["type"] = "Point"
	gjson["coordinates"] = []float64{center[1], center[0]}

	geometry := make(map[string]interface{})
	geometry["$geometry"] = gjson
	// geometry["$maxDistance"] = 40010000

	gw := make(map[string]interface{})
	gw["$nearSphere"] = geometry

	qi := make(map[string]interface{})
	qi["geo.geometry"] = gw

	return mdb.Query(qi)[0:count]
}

func (mdb *MongoDatabase) sanitize(string_id string) string {
	return strings.Replace(strings.Replace(string_id, "ObjectIdHex(\"", "", 1), "\")", "", 1)
}
