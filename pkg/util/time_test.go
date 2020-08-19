package util

import (
	"testing"
	"time"

	"github.com/matryer/is"
)

func TestParseDuration(t *testing.T) {
	is := is.New(t)

	_, err := ParseDuration("15x")
	is.True(err != nil)

	dur, err := ParseDuration("80s")
	is.NoErr(err)
	is.Equal(dur, time.Second*80)
}
