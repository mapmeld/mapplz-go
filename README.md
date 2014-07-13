# MapPLZ-Go

[MapPLZ](http://mapplz.com) is a framework to make mapping quick and easy in
your favorite language.

<img src="https://raw.githubusercontent.com/mapmeld/mapplz-go/master/logo.jpg" width="140"/>

## Getting started

MapPLZ consumes many many types of geodata. It can process data for a script or dump
it into a database.

Go does not support method overloading, so you need to choose the right function to
submit your data.

Here's how you can add some data:

```
// add points
mapplz.Add_Lat_Lng(40, -70)
mapplz.Add_Lng_Lat(-70, 40)
mapplz.Add_LatLng([40, -70])

// add GeoJSON
mapplz.Add_Geojson_Feature_Str(`{ "type": "Feature", "geometry": { "type": "Point", "coordinates": [-70, 40] }}`)
```

## License

Free BSD License
