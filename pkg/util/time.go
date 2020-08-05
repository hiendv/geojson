package util

import (
	"time"
)

// ParseDuration parses a duration string.
func ParseDuration(str string) (duration time.Duration, err error) {
	duration, err = time.ParseDuration(str)
	if err != nil {
		return
	}

	return
}
