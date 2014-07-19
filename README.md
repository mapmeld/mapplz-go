# MapPLZ-Go

[MapPLZ](http://mapplz.com) is a framework to make mapping quick and easy in
your favorite language.

<img src="https://raw.githubusercontent.com/mapmeld/mapplz-go/master/logo.jpg" width="140"/>

## Getting started

MapPLZ consumes many many types of geodata. It can process data for a script or dump
it into a database.

Go does not support method overloading, so you need to name the right function for
your data. Parameters are separated by an underscore, so sending (lat, lng) would be
Add_Lat_Lng, sending (lng, lat) would be Add_Lng_Lat, and sending a single parameter
{lat, lng} is Add_LatLng.

Adding some data:

```
mapstore := mapplz.NewMapPLZ()

// add points
mapstore.Add_Lat_Lng(40, -70)
mapstore.Add_Lng_Lat(-70, 40)
mapstore.Add_LatLng([40, -70])

// add lines
mapstore.Add_LatLngPath([[40, -70], [50, 20]])

// add polygons
mapstore.Add_LatLngPoly([[40, -70], [50, 20], [40, 40], [40, -70]])

// GeoJSON strings (points, lines, and polygons)
mapstore.Add_Geojson_Feature(`{ "type": "Feature", "geometry": { "type": "Point", "coordinates": [-70, 40] }}`)
```

Each feature is added to the mapstore and returned as a MapItem

```
pt := mapstore.Add_Lat_Lng(40, -70)
line := mapstore.Add_LatLngPath([[40, -70], [50, 20]])

len(mapstore.MapItems) == 2

pt.Lat() == 40
pt.ToGeoJson() == `{ "type": "Feature", "geometry": { "type": "Point", "coordinates": [-70, 40] }}`

line.Type() == "line"
line.Path() == [[[40, -70], [50, 20]]]

```

## Packages

* <a href="https://github.com/kellydunn/golang-geo">golang-geo</a> from Kelly Dunn (MIT license)
* <a href="https://github.com/kpawlik/geojson">geojson</a> from Kris Pawlik (MIT license)

## License

Free BSD License
