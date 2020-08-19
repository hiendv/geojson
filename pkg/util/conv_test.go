package util

import (
	"testing"

	"github.com/matryer/is"
)

func TestInt64FromString(t *testing.T) {
	is := is.New(t)

	_, err := Int64FromString("a")
	is.True(err != nil)

	x, err := Int64FromString("12")
	is.NoErr(err)
	is.Equal(x, int64(12))
}
