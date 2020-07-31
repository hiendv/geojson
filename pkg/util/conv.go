package util

import (
	"strconv"
)

func Int64FromString(s string) (int64, error) {
	int, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, err
	}

	return int64(int), nil
}
