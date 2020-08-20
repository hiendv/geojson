package osm

import (
	"context"
	"errors"
	"os"
)

func validateOut(path string) error {
	_, err := os.Stat(path)
	return err
}

// FindSubAreas looks for outputs of a sub-area.
func FindSubAreas(ctx context.Context, id int64) (string, error) {
	path, ok := filePath(ctx, id)
	if !ok {
		return "", errors.New("invalid directory")
	}

	err := VerifyOutput(ctx, path)
	if err != nil {
		return "", err
	}

	return path, nil
}

// VerifyOutput makes sure that output files of a sub-area exist.
func VerifyOutput(ctx context.Context, path string) error {
	_, err := os.Stat(path)
	if err != nil {
		return errors.New("invalid path")
	}

	return nil
}
