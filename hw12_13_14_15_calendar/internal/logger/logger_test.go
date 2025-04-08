package logger

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	t.Run("test logger level", func(t *testing.T) {
		log := New(ctx, "debug", nil, true)
		require.Equal(t, DebugLevel, log.level)
	})

	t.Run("test Info method", func(t *testing.T) {
		var buf bytes.Buffer
		log := New(ctx, "info", nil, true).Output(&buf)
		log.Info("info message")
		require.Contains(t, buf.String(), "info message")
	})

	t.Run("test Error method", func(t *testing.T) {
		var buf bytes.Buffer
		log := New(ctx, "error", nil, true).Output(&buf)
		log.Error("error message")
		require.Contains(t, buf.String(), "error: error message")
	})

	t.Run("test Output method", func(t *testing.T) {
		var buf1, buf2 bytes.Buffer
		log := New(ctx, "debug", nil, true).Output(&buf1)
		log = log.Output(&buf2)
		log.Info("info message")
		require.Contains(t, buf2.String(), "info: info message")
	})
}
