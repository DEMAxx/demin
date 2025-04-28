package logger

import (
	"context"
	"fmt"
	"io"
	"runtime/debug"
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
		return l.Write(append([]byte("\n info: "), p...))
	case ErrorLevel:
		return l.Write(append([]byte("error: "), p...))
	case WarnLevel:
		return l.Write(append([]byte("warn: "), p...))
	case TraceLevel:
		return l.Write(append([]byte("trace: "), p...))
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
	case TraceLevel:
		return "trace"
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

func New(ctx context.Context, level string, sampler Sampler, stack bool) *Logger {
	var logs Logger
	logs.ctx = ctx
	logs.sampler = sampler
	logs.stack = stack

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

func NewWriter(ctx context.Context, w io.Writer, sampler Sampler, stack bool) Logger {
	if w == nil {
		w = io.Discard
	}
	lw, ok := w.(LevelWriter)
	if !ok {
		lw = LevelWriterAdapter{w}
	}
	return Logger{w: lw, level: TraceLevel, ctx: ctx, sampler: sampler, stack: stack}
}

func (l Logger) Info(msg string) {
	if l.sampler != nil && !l.sampler.Sample(InfoLevel) {
		return
	}
	_, err := l.w.WriteLevel(InfoLevel, append(l.context, []byte(msg)...))
	if err != nil {
		return
	}
}

func (l Logger) Error(msg string) {
	if l.sampler != nil && !l.sampler.Sample(ErrorLevel) {
		return
	}

	var message []byte

	message = append(message, []byte(fmt.Sprintf("error: %s", msg))...)
	if l.stack {
		message = append(message, []byte("\n"+string(debug.Stack()))...)
	}
	_, err := l.w.WriteLevel(ErrorLevel, message)
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
