# spatial

![](https://travis-ci.org/dewski/spatial.svg?branch=master)

Add simple types to use PostGIS with any database driver in golang.

Types supported

- geometry(Geometry,4326)

## Usage with database

First define your struct:

```go
package main

type TripEvent struct {
  Geo spatial.Point `db:"geo" json:"geo"`
}

func main() {
  path := []spatial.Point{}
  err := database.DB.
  	Select("geo").
  	From("trip_events").
  	QueryStructs(&path)

  if err != nil {
    panic(err)
  }

  // Encode the path with .000000 level of precision
  polyline := spatial.Encode(path, 6)
  fmt.Println(polyline) // _p~iF~ps|U_ulLnnqC_mqNvxq`@

  points := spatial.Decode(polyline, 6)
  fmt.Println(polyline) // make(map[]spatial.Point, 2)
}
```

If your users aren't going to be zooming you can ignore the level 6 precision
and go with 5 to get the best compression.


## Usage without database

First define your struct:

```go
package main

func main() {
  path := []spatial.Point{
    spatial.Point{
      Lat: 38.889803,
      Lng: -77.009114,
    },
    spatial.Point{
      Lat: 38.889810,
      Lng: -77.009124,
    },
  }

  // Encode the path with .000000 level of precision
  polyline := spatial.Encode(path, 6)
  fmt.Println(polyline) // _p~iF~ps|U_ulLnnqC_mqNvxq`@

  points := spatial.Decode(polyline, 6)
  fmt.Println(polyline) // make(map[]spatial.Point, 2)
}
```
