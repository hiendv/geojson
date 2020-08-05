package osm

import (
	"context"
	"errors"
	"os"
)

func FindSubAreas(ctx context.Context, id int64) (string, error) {
	path, ok := filePath(ctx, id)
	if !ok {
		return "", errors.New("invalid directory")
	}
	_, err := os.Stat(path)
	if err != nil {
		return "", errors.New("invalid path")
	}

	return path, nil
}
