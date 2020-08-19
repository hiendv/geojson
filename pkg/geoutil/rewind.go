package geoutil

import (
	"errors"

	"github.com/hiendv/geojson/pkg/util"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
)

const geometryMultiPolygon = "MultiPolygon"
const geometryPolygon = "Polygon"

func RewindFeatureCollection(fc *geojson.FeatureCollection, outer bool) error {
	if fc == nil {
		return errors.New("invalid feature collection")
	}

	for _, feature := range fc.Features {
		err := RewindFeature(feature, outer)
		if err != nil {
			return err
		}
	}

	return nil
}

func RewindFeature(f *geojson.Feature, outer bool) error {
	if f == nil {
		return errors.New("invalid feature")
	}

	geometry, err := RewindGeometry(f.Geometry, outer)
	if err != nil {
		return err
	}

	f.Geometry = geometry

	return nil
}

func RewindGeometry(g orb.Geometry, outer bool) (orb.Geometry, error) {
	if g == nil {
		return nil, errors.New("invalid geometry")
	}

	if g.GeoJSONType() == geometryPolygon {
		mp, ok := g.(orb.Polygon)
		if !ok {
			return nil, errors.New("invalid Polygon")
		}

		mp = RewindRings(mp, outer)
		return mp, nil
	}

	if g.GeoJSONType() == geometryMultiPolygon {
		mp, ok := g.(orb.MultiPolygon)
		if !ok {
			return nil, errors.New("invalid MultiPolygon")
		}

		for i, p := range mp {
			mp[i] = RewindRings(p, outer)
		}

		return mp, nil
	}

	return g, errors.New("geometry type not supported")
}

func RewindRings(rings []orb.Ring, outer bool) []orb.Ring {
	if len(rings) == 0 {
		return rings
	}

	RewindRing(rings[0], outer)
	for i := 1; i < len(rings); i++ {
		RewindRing(rings[i], !outer)
	}

	return rings
}

func RewindRing(ring orb.Ring, cw bool) {
	// Shoelace formula: https://mathworld.wolfram.com/PolygonArea.html
	var area float64 = 0
	len := len(ring)
	for i, j := 0, len-1; i < len; i, j = i+1, i+1 {
		area += (ring[i][0] - ring[j][0]) * (ring[j][1] + ring[i][1])
	}

	if area >= 0 != cw {
		util.ReverseAny(ring)
		return
	}

	return
}
