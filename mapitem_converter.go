package mapplz

import (
	"encoding/json"
)

func Convert_Lat_Lng(lat float64, lng float64, db MapDatabase) MapItem {
	return NewMapItemPoint(lat, lng, db)
}

func ConvertPath(path [][]float64, db MapDatabase) MapItem {
	return NewMapItemLine(path, db)
}

func ConvertPoly(poly [][]float64, db MapDatabase) MapItem {
	return NewMapItemPoly(poly, db)
}

func ConvertGeojsonFeature(geojson string, db MapDatabase) MapItem {
	var mip MapItem

	// GeoJSON parsing based on http://stackoverflow.com/a/15728702
	var geojsonData GeojsonFeature
	err := json.Unmarshal([]byte(geojson), &geojsonData)
	if err != nil {
		panic("Failed to parse JSON string")
	}

	switch geojsonData.Geometry.Type {
	case "Point":
		json.Unmarshal(geojsonData.Geometry.Coordinates, &geojsonData.Geometry.Point.Coordinates)
		mip = Convert_Lat_Lng(geojsonData.Geometry.Point.Coordinates[1], geojsonData.Geometry.Point.Coordinates[0], db)
	case "LineString":
		json.Unmarshal(geojsonData.Geometry.Coordinates, &geojsonData.Geometry.Line.Coordinates)
		path := geojsonData.Geometry.Line.Coordinates
		for i := 0; i < len(path); i++ {
			lat := path[i][1]
			lng := path[i][0]
			path[i][0] = lat
			path[i][1] = lng
		}
		mip = ConvertPath(path, db)
	case "Polygon":
		json.Unmarshal(geojsonData.Geometry.Coordinates, &geojsonData.Geometry.Polygon.Coordinates)
		path := geojsonData.Geometry.Polygon.Coordinates[0]
		for i := 0; i < len(path); i++ {
			lat := path[i][1]
			lng := path[i][0]
			path[i][0] = lat
			path[i][1] = lng
		}
		mip = ConvertPoly(path, db)
	default:
		panic("Unsupported GeoJSON Feature Type")
	}

	mip.SetProperties(geojsonData.Properties)

	return mip
}
