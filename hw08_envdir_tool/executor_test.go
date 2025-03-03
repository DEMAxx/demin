package main

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRunCmd(t *testing.T) {
	t.Run("no command provided", func(t *testing.T) {
		cmd := []string{}
		env := Environment{}
		returnCode := RunCmd(cmd, env)
		require.Equal(t, 1, returnCode)
	})

	t.Run("command with environment variables", func(t *testing.T) {
		cmd := []string{"env"}
		env := Environment{
			"FOO":   {Value: "foo", NeedRemove: false},
			"BAR":   {Value: "bar", NeedRemove: false},
			"UNSET": {Value: "", NeedRemove: true},
		}
		returnCode := RunCmd(cmd, env)
		require.Equal(t, 0, returnCode)
	})

	t.Run("command with empty environment variable", func(t *testing.T) {
		cmd := []string{"env"}
		env := Environment{
			"EMPTY": {Value: "", NeedRemove: false},
		}
		returnCode := RunCmd(cmd, env)
		require.Equal(t, 0, returnCode)
	})

	t.Run("command with environment variable to remove", func(t *testing.T) {
		cmd := []string{"env"}
		env := Environment{
			"REMOVE_ME": {Value: "", NeedRemove: true},
		}
		returnCode := RunCmd(cmd, env)
		require.Equal(t, 0, returnCode)
	})
}
