package logger

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	t.Run("test logger level", func(t *testing.T) {
		log := New("debug")
		require.Equal(t, DebugLevel, log.level)
	})

	t.Run("test Info method", func(t *testing.T) {
		var buf bytes.Buffer
		log := New("info").Output(&buf)
		log.Info("info message")
		require.Contains(t, buf.String(), "info message")
	})

	t.Run("test Error method", func(t *testing.T) {
		var buf bytes.Buffer
		log := New("error").Output(&buf)
		log.Error("error message")
		require.Contains(t, buf.String(), "error: error message")
	})

	t.Run("test Output method", func(t *testing.T) {
		var buf1, buf2 bytes.Buffer
		log := New("debug").Output(&buf1)
		log = log.Output(&buf2)
		log.Info("info message")
		require.Contains(t, buf2.String(), "info: info message")
	})
}
