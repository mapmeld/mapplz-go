# MapPLZ-Go

[MapPLZ](http://mapplz.com) is a framework to make mapping quick and easy in
your favorite language.

<img src="https://raw.githubusercontent.com/mapmeld/mapplz-go/master/logo.jpg" width="140"/>

## Getting started

MapPLZ consumes many many types of geodata. It can process data for a script or dump
it into a database.

Go does not support method overloading, so you need to use Add, Add2, or Add3 depending
on the number of parameters.

You can also call specific Add functions. Function names are based on their parameters
separated by an underscore, so send (lat, lng) to Add_Lat_Lng, send (lng, lat)
to Add_Lng_Lat, and send an array [lat, lng] to Add_LatLng.

Adding some data:

```go
mapstore := mapplz.NewMapPLZ()

// add points
mapstore.Add_Lat_Lng(40, -70)
mapstore.Add_Lng_Lat(-70, 40)
mapstore.Add_LatLng( []float64{ 40, -70 } )

// add lines
mapstore.Add_LatLngPath([][]float64{{40, -70}, {50, 20}})

// add polygons
mapstore.Add_LatLngPoly([][]float64{{40, -70}, {50, 20}, {40, 40}, {40, -70}})

// GeoJSON strings (points, lines, and polygons)
mapstore.Add(`{ "type": "Feature", "geometry": { "type": "Point", "coordinates": [-70, 40] }}`)

// add properties
mapstore.Add2([]float64{40, -70}, `{ "color": "red" }`)
mapstore.Add(`{ "type": "Feature", "geometry": { "type": "LineString", "coordinates": [[-70, 40], [-80, 50]] }, "properties": { "color": "#0f0" }}`)

props := make(map[string]interface{})
props["lat"] = 40
props["lng"] = -70
props["property"] = "test"
mapstore.Add(props)
```

Each feature is added to the mapstore and returned as a MapItem

```go
pt := mapstore.Add_Lat_Lng(40, -70)
pt.SetJsonProperties(`{ "color": "#00f" }`)
pt.SetProperties(map[string]interface{})
pt.Properties()["color"]
pt.Lat()
pt.ToGeoJson() == `{ "type": "Feature", "geometry": { "type": "Point", "coordinates": [-70, 40] }, "properties": { "color": "#00f" }}`
pt.Delete()

line.Type() == "line"
line.Path() == [[[40, -70], [50, 20]]]
line.Properties()["color"]

mapstore.ToGeoJson() // export all as GeoJSON
mapstore.Count("") == 2  // MapItems count ("" selects all)
mapstore.Query("") // array of all MapItems ("" selects all)

// supported with PostGIS
mapstore.Where("color = 'red'") // enter a SQL WHERE clause with one property
mapstore.Count("color = 'white' OR color = 'blue'")
item.Delete()

// supported with MongoDB or without DB
query = make(map[string]interface{})
query["color"] = "red"
mapstore.Query(query)
mapstore.Count(query)
item.Delete()
```

You can make some geospatial queries:

```go
// all MapItems inside this area
mapstore.Within(`{ "type": "Feature", "geometry": { "type": "Polygon", "coordinates": [[[ ... ]]] } }`)
mapstore.Within([][]float64{{40, -70}, {40, -110}, {60, -90}, {40, -70}})

// MapItems nearest to a point
mapstore.Near([]float64{40, -70}, 2)
mapstore.Near(`{ "type": "Feature", "geometry": { "type": "Point", "coordinates": [-70, 40] } }`, 2)

// point in polygon?
pt.Within([][]float64{{40, -70}, {40, -110}, {60, -90}, {40, -70}})

// distance between points?
pt.DistanceFrom([]float64{lat, lng})
```

## Instant Interactive Maps

Get the HTML and JavaScript code to display your data on a map:

```go
// HTML embed code
mapstore.EmbedHtml()

// full HTML page (including head, body, stylesheet)
mapstore.RenderHtml()
```

## Databases

MapPLZ can set up geospatial data with PostGIS or MongoDB, and take the complexities out of your hands.

Here's how you can connect:

#### PostGIS
```go
// install PostGIS and create a 'mapplz' table before use
// here's my schema:
// CREATE TABLE mapplz (id SERIAL PRIMARY KEY, properties JSON, geom public.geometry)

mapstore := NewMapPLZ()
db, _ := sql.Open("postgres", "user=USER dbname=DB sslmode=SSLMODE")
mapstore.Database = NewPostGISDB(db)
```

#### MongoDB

```go
// install MongoDB and create a db and collection
mapstore := NewMapPLZ()
session, err := mgo.Dial("localhost")
defer session.Close()
collection := session.DB("sample").C("mapplz")
// put a geospatial index on 'geo.geometry'
geoindex := mgo.Index{Key: []string{"$2dsphere:geo.geometry"}, Bits: 26}
collection.EnsureIndex(geoindex)
mapstore.Database = NewMongoDatabase(collection)
```

## Packages

* <a href="https://github.com/kellydunn/golang-geo">golang-geo</a> from Kelly Dunn (MIT license)
* <a href="https://github.com/mapmeld/geojson-bson">geojson-bson</a> based on <a href="https://github.com/kpawlik/geojson">geojson</a> from Kris Pawlik (MIT license)
* <a href="https://github.com/lib/pq">pq</a> (MIT license)
* <a href="http://gopkg.in/mgo.v2">mgo</a> (Simplified BSD license)

With additional input from

* <a href="https://github.com/akavel/polyclip-go">polyclip-go</a> (MIT license) from Mateusz Czapliński

## License

Free BSD License
