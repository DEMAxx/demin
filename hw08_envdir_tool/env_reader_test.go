package main

import (
	"errors"
	"fmt"
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
		dir := "./testdata/env"

		env, err := ReadDir(dir)

		if err != nil {
			println("error:", err.Error())
		}

		for key, val := range env {
			fmt.Printf("test %s : %s\n", key, val.Value)

		}

		println(env)

		//todo finish test to check is env have empty values
	})
}
