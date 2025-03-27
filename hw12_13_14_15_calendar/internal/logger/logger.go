package logger

import (
	"context"
	"fmt"
	"io"
)

type Level int8

type LevelWriter interface {
	io.Writer
	WriteLevel(level Level, p []byte) (n int, err error)
}

type LevelWriterAdapter struct {
	io.Writer
}

func (l LevelWriterAdapter) WriteLevel(level Level, p []byte) (n int, err error) {
	switch level {
	case DebugLevel:
		return l.Write(append([]byte("debug: "), p...))
	case InfoLevel:
		return l.Write(append([]byte("info: "), p...))
	case ErrorLevel:
		return l.Write(append([]byte("error: "), p...))
	default:
		return l.Write(p)
	}
}

type Sampler interface {
	Sample(lvl Level) bool
}

const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	TraceLevel Level = -1
)

func (l Level) String() string {
	switch l {
	case DebugLevel:
		return "debug"
	case InfoLevel:
		return "info"
	case WarnLevel:
		return "warn"
	case ErrorLevel:
		return "error"
	default:
		return ""
	}
}

type Logger struct {
	w       LevelWriter
	level   Level
	sampler Sampler
	context []byte
	stack   bool
	ctx     context.Context
}

func New(level string) *Logger {

	var logs Logger

	switch level {
	case "trace":
		logs.Level(TraceLevel)
	case "debug":
		logs = logs.Level(DebugLevel)
	case "warn":
		logs = logs.Level(WarnLevel)
	case "error":
		logs = logs.Level(ErrorLevel)
	default:
		logs = logs.Level(InfoLevel)
	}

	return &logs
}

func NewWriter(w io.Writer) Logger {
	if w == nil {
		w = io.Discard
	}
	lw, ok := w.(LevelWriter)
	if !ok {
		lw = LevelWriterAdapter{w}
	}
	return Logger{w: lw, level: TraceLevel}
}

func (l Logger) Info(msg string) {
	_, err := l.w.WriteLevel(InfoLevel, []byte(msg))

	if err != nil {
		return
	}
}

func (l Logger) Error(msg string) {
	_, err := l.w.WriteLevel(ErrorLevel, []byte(fmt.Sprintf("error: %s", msg)))

	if err != nil {
		return
	}
}

func (l Logger) Level(lvl Level) Logger {
	l.level = lvl
	return l
}

func (l *Logger) Output(w io.Writer) *Logger {
	lw, ok := w.(LevelWriter)
	if !ok {
		lw = LevelWriterAdapter{w}
	}
	l.w = lw
	return l
}
