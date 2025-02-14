package main

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	t.Run("error file not found", func(t *testing.T) {
		err := Copy("/test/error.txt", "out.txt", 0, 0)
		require.Truef(t, errors.Is(err, ErrUnsupportedFile), "actual err - %v", err)
	})

	t.Run("copy limit work", func(t *testing.T) {
		err := Copy("./testdata/input.txt", "out.txt", 0, 1000)

		require.Truef(t, errors.Is(err, nil), "actual err - %v", err)
	})

	t.Run("copy offset work", func(t *testing.T) {
		err := Copy("./testdata/input.txt", "out.txt", 1000, 0)

		require.Truef(t, errors.Is(err, nil), "actual err - %v", err)
	})

	t.Run("error offset is more then filesize", func(t *testing.T) {
		err := Copy("./testdata/input.txt", "out.txt", 10000000, 0)

		require.Truef(t, errors.Is(err, ErrOffsetExceedsFileSize), "actual err - %v", err)
	})
}
