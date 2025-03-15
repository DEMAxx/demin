package logger

import (
	"context"
	"fmt"
	"io"
	"os"
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
	return l.Write(p)
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

	logs = logs.Output(os.Stderr)

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
	fmt.Println(msg)
}

func (l Logger) Error(msg string) {
	_ = fmt.Errorf("error: %s", msg)
}

func (l Logger) Level(lvl Level) Logger {
	l.level = lvl
	return l
}

func (l Logger) Output(w io.Writer) Logger {
	l2 := NewWriter(w)
	l2.level = l.level
	l2.sampler = l.sampler
	l2.stack = l.stack

	if l.context != nil {
		l2.context = make([]byte, len(l.context), cap(l.context))
		copy(l2.context, l.context)
	}
	return l2
}
