package geoutil

import (
	"encoding/json"
	"testing"

	"github.com/matryer/is"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
)

/*
   B +-------+ C
     |       |
     |       |
     |       |
   A +-------+ D
*/
func TestRewindRingEmpty(t *testing.T) {
	is := is.New(t)
	ring := orb.Ring{}

	RewindRing(ring, true)
	is.Equal(ring, orb.Ring{})

	RewindRing(ring, false)
	is.Equal(ring, orb.Ring{})
}

func TestRewindRingClockwise(t *testing.T) {
	is := is.New(t)
	ring := orb.Ring{
		orb.Point{1, 1}, // A
		orb.Point{1, 2}, // B
		orb.Point{2, 2}, // C
		orb.Point{2, 1}, // D
		orb.Point{1, 1}, // A
	}

	RewindRing(ring, true) // true = CW
	is.Equal(ring, orb.Ring{
		orb.Point{1, 1}, // A
		orb.Point{1, 2}, // B
		orb.Point{2, 2}, // C
		orb.Point{2, 1}, // D
		orb.Point{1, 1}, // A
	})
}

func TestRewindRingCounterClockwise(t *testing.T) {
	is := is.New(t)
	ring := orb.Ring{
		orb.Point{1, 1}, // A
		orb.Point{1, 2}, // B
		orb.Point{2, 2}, // C
		orb.Point{2, 1}, // D
		orb.Point{1, 1}, // A
	}

	RewindRing(ring, false) // false = CCW
	is.Equal(ring, orb.Ring{
		orb.Point{1, 1}, // A
		orb.Point{2, 1}, // D
		orb.Point{2, 2}, // C
		orb.Point{1, 2}, // B
		orb.Point{1, 1}, // A
	})
}

/*
   F +---------------+ G
     |               |
     | B +-------+ C |
     |   |       |   |
     |   |       |   |
     |   |       |   |
     | A +-------+ D |
     |               |
   E +---------------+ H
*/
func TestRewindRingsEmpty(t *testing.T) {
	is := is.New(t)
	ring := []orb.Ring{}

	RewindRings(ring, true)
	is.Equal(ring, []orb.Ring{})

	RewindRings(ring, false)
	is.Equal(ring, []orb.Ring{})
}

func TestRewindRingsInverseRFC7946(t *testing.T) {
	/*
		https://tools.ietf.org/html/rfc7946#section-3.1.6
		A linear ring MUST follow the right-hand rule with respect to the area it bounds, i.e., exterior rings are counterclockwise, and holes are clockwise.
	*/

	is := is.New(t)
	rings := []orb.Ring{
		orb.Ring{
			orb.Point{1, 1}, // A
			orb.Point{1, 2}, // B
			orb.Point{2, 2}, // C
			orb.Point{2, 1}, // D
			orb.Point{1, 1}, // A
		},
		orb.Ring{
			orb.Point{0, 0}, // E
			orb.Point{0, 3}, // F
			orb.Point{3, 3}, // G
			orb.Point{3, 0}, // H
			orb.Point{0, 0}, // E
		},
	}

	RewindRings(rings, false) // false = CCW
	is.Equal(rings, []orb.Ring{
		orb.Ring{
			orb.Point{1, 1}, // A
			orb.Point{2, 1}, // D
			orb.Point{2, 2}, // C
			orb.Point{1, 2}, // B
			orb.Point{1, 1}, // A
		},
		orb.Ring{
			orb.Point{0, 0}, // E
			orb.Point{0, 3}, // F
			orb.Point{3, 3}, // G
			orb.Point{3, 0}, // H
			orb.Point{0, 0}, // E
		},
	})
}

func TestRewindGeoEmpty(t *testing.T) {
	is := is.New(t)

	err := RewindGeometry(nil, false)
	is.True(err != nil)
}

func TestRewindGeoInvalid(t *testing.T) {
	is := is.New(t)

	err := RewindGeometry(&invalidPolygon{}, false)
	is.True(err != nil)

	err = RewindGeometry(&invalidMultiPolygon{}, false)
	is.True(err != nil)

	err = RewindGeometry(orb.MultiPoint(nil), false)
	is.True(err != nil)
}

func TestRewindPolygon(t *testing.T) {
	is := is.New(t)
	var geo, rewinded orb.Geometry

	geo = orb.Polygon([]orb.Ring{
		orb.Ring{
			orb.Point{1, 1}, // A
			orb.Point{1, 2}, // B
			orb.Point{2, 2}, // C
			orb.Point{2, 1}, // D
			orb.Point{1, 1}, // A
		},
		orb.Ring{
			orb.Point{0, 0}, // E
			orb.Point{0, 3}, // F
			orb.Point{3, 3}, // G
			orb.Point{3, 0}, // H
			orb.Point{0, 0}, // E
		},
	})

	rewinded = orb.Polygon([]orb.Ring{
		orb.Ring{
			orb.Point{1, 1}, // A
			orb.Point{2, 1}, // D
			orb.Point{2, 2}, // C
			orb.Point{1, 2}, // B
			orb.Point{1, 1}, // A
		},
		orb.Ring{
			orb.Point{0, 0}, // E
			orb.Point{0, 3}, // F
			orb.Point{3, 3}, // G
			orb.Point{3, 0}, // H
			orb.Point{0, 0}, // E
		},
	})

	err := RewindGeometry(geo, false) // false = CCW
	is.NoErr(err)
	is.Equal(geo, rewinded)
}

func TestRewindMultiPolygon(t *testing.T) {
	is := is.New(t)
	var geo, rewinded orb.Geometry

	geo = orb.MultiPolygon([]orb.Polygon{[]orb.Ring{
		orb.Ring{
			orb.Point{1, 1}, // A
			orb.Point{1, 2}, // B
			orb.Point{2, 2}, // C
			orb.Point{2, 1}, // D
			orb.Point{1, 1}, // A
		},
		orb.Ring{
			orb.Point{0, 0}, // E
			orb.Point{0, 3}, // F
			orb.Point{3, 3}, // G
			orb.Point{3, 0}, // H
			orb.Point{0, 0}, // E
		},
	}})

	rewinded = orb.MultiPolygon([]orb.Polygon{[]orb.Ring{
		orb.Ring{
			orb.Point{1, 1}, // A
			orb.Point{2, 1}, // D
			orb.Point{2, 2}, // C
			orb.Point{1, 2}, // B
			orb.Point{1, 1}, // A
		},
		orb.Ring{
			orb.Point{0, 0}, // E
			orb.Point{0, 3}, // F
			orb.Point{3, 3}, // G
			orb.Point{3, 0}, // H
			orb.Point{0, 0}, // E
		},
	}})

	err := RewindGeometry(geo, false) // false = CCW
	is.NoErr(err)
	is.Equal(geo, rewinded)
}

func TestRewindFeatureInvalid(t *testing.T) {
	is := is.New(t)

	err := RewindFeature(nil, false)
	is.True(err != nil)

	feature := geojson.NewFeature(nil)
	err = RewindFeature(feature, false)
	is.True(err != nil)
}

func TestRewindFeature(t *testing.T) {
	is := is.New(t)

	feature := geojson.NewFeature(orb.Polygon([]orb.Ring{
		orb.Ring{
			orb.Point{1, 1}, // A
			orb.Point{1, 2}, // B
			orb.Point{2, 2}, // C
			orb.Point{2, 1}, // D
			orb.Point{1, 1}, // A
		},
		orb.Ring{
			orb.Point{0, 0}, // E
			orb.Point{0, 3}, // F
			orb.Point{3, 3}, // G
			orb.Point{3, 0}, // H
			orb.Point{0, 0}, // E
		},
	}))

	rewinded := geojson.NewFeature(orb.Polygon([]orb.Ring{
		orb.Ring{
			orb.Point{1, 1}, // A
			orb.Point{2, 1}, // D
			orb.Point{2, 2}, // C
			orb.Point{1, 2}, // B
			orb.Point{1, 1}, // A
		},
		orb.Ring{
			orb.Point{0, 0}, // E
			orb.Point{0, 3}, // F
			orb.Point{3, 3}, // G
			orb.Point{3, 0}, // H
			orb.Point{0, 0}, // E
		},
	}))

	err := RewindFeature(feature, false) // false = CCW
	is.NoErr(err)
	is.Equal(feature, rewinded)
}

func TestRewindFeatureCollectionInvalid(t *testing.T) {
	is := is.New(t)

	err := RewindFeatureCollection(nil, false)
	is.True(err != nil)

	fc, err := geojson.UnmarshalFeatureCollection([]byte(`
		{
			"type": "FeatureCollection",
			"features": [
				{
					"type": "Feature",
					"geometry": {
						"type": "MultiPolygon"
					}
				}
			]
		}
	`))
	is.NoErr(err)

	fc.Features = []*geojson.Feature{geojson.NewFeature(nil)}
	err = RewindFeatureCollection(fc, false)
	is.True(err != nil)
}

func TestRewindFeatureCollection(t *testing.T) {
	is := is.New(t)

	fc, err := geojson.UnmarshalFeatureCollection([]byte(`
		{
			"type": "FeatureCollection",
			"features": [
				{
					"type": "Feature",
					"geometry": {
						"type": "MultiPolygon",
						"coordinates": [
							[
								[
									[1, 1],
									[1, 2],
									[2, 2],
									[2, 1],
									[1, 1]
								],
								[
									[0, 0],
									[0, 3],
									[3, 3],
									[3, 0],
									[0, 0]
								]
							]
						]
					},
					"properties": null
				}
			]
		}
	`))
	is.NoErr(err)

	rewinded := []byte(`{"type":"FeatureCollection","features":[{"type":"Feature","geometry":{"type":"MultiPolygon","coordinates":[[[[1,1],[2,1],[2,2],[1,2],[1,1]],[[0,0],[0,3],[3,3],[3,0],[0,0]]]]},"properties":null}]}`)

	err = RewindFeatureCollection(fc, false) // false = CCW
	is.NoErr(err)

	x, err := json.Marshal(fc)
	is.NoErr(err)
	is.Equal(x, rewinded)
}
