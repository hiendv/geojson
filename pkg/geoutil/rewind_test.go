package geoutil

import (
	"testing"

	"github.com/matryer/is"
	"github.com/paulmach/orb"
)

/*
   B +-------+ C
     |       |
     |       |
     |       |
   A +-------+ D
*/

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
