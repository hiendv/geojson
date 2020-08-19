package geoutil

import (
	"github.com/paulmach/orb"
)

type invalidPolygon struct {
	orb.MultiPolygon
}

func (p *invalidPolygon) GeoJSONType() string {
	return geometryPolygon
}

func (p *invalidPolygon) Dimensions() int {
	return 2
}

func (p *invalidPolygon) Bound() orb.Bound {
	return orb.Bound{Min: orb.Point{1, 1}, Max: orb.Point{-1, -1}}
}

func (p *invalidPolygon) private() {
}

type invalidMultiPolygon struct {
	orb.Polygon
}

func (p *invalidMultiPolygon) GeoJSONType() string {
	return geometryMultiPolygon
}

func (p *invalidMultiPolygon) Dimensions() int {
	return 2
}

func (p *invalidMultiPolygon) Bound() orb.Bound {
	return orb.Bound{Min: orb.Point{1, 1}, Max: orb.Point{-1, -1}}
}

func (p *invalidMultiPolygon) private() {
}
