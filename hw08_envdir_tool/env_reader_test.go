package main

import (
	"errors"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestReadDir(t *testing.T) {

	t.Run("wrong dir", func(t *testing.T) {
		dir := "./test/env"

		_, err := ReadDir(dir)

		require.Truef(t, errors.Is(err, ErrInvalidDir), "actual err - %v", err)
	})

	t.Run("no empty values", func(t *testing.T) {
		dir := "./testdata/env2"

		env, err := ReadDir(dir)

		if err != nil {
			println("error:", err.Error())
		}

		for key, val := range env {
			require.Truef(t, val.Value != "", "error empty value: %v", key)
		}
	})
}
