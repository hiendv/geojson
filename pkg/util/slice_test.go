package util

import (
	"testing"

	"github.com/matryer/is"
)

func TestReverseAny(t *testing.T) {
	is := is.New(t)

	ints := []int{1, 2, 3, 4}
	strs := []string{"a", "b", "c", "d"}

	ReverseAny(ints)
	ReverseAny(strs)

	is.Equal(ints, []int{4, 3, 2, 1})
	is.Equal(strs, []string{"d", "c", "b", "a"})
}
