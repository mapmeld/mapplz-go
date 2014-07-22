# Wishlist (Go-specific)

## Go Lang

* Use better function names
* Ask community for feedback
* Separate out testing of different types
* Handle and pass along a massive amount of errors
* Is there a way to install all dependencies?

## Databases

* Separate out geo ETL into a static method, so it can be called from a MapDatabase
* Refactor geo ETL so it can be called without adding to MapItems
* Save IDs
* Add Save() to MapItem

## Formats

* Return []MapItem from GeoJSON FeatureCollection, remove old types
* Make a real ToWKT() on Line and Poly MapItems
* Always return an array of MapItems (or maybe [0] if there's only one)

## No DB

* Query by property
* Spatial queries

## PostGIS

* Load geo results without putting them into a MapPLZ instance
* Test saving and loading properties
* Count
* Query by property
* Spatial queries

## MongoDB

* Document bzr install, get working on Travis CI
* Write

## Spatialite?
