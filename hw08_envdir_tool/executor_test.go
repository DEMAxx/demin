package main

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRunCmd(t *testing.T) {
	t.Run("executed", func(t *testing.T) {
		cmd := []string{
			"test1",
		}

		env, err := ReadDir("./testdata/env")

		require.Truef(t, err == nil, "error reading dir: %v", err)

		rCmd := RunCmd(cmd, env)

		require.True(t, rCmd == 1)
	})
}
