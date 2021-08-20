package util

import (
	"io"
	"time"
)

type Logger interface {
	Printf(format string, v ...interface{})
	Println(v ...interface{})
}

type NopLogger struct{}

func (l *NopLogger) Printf(format string, v ...interface{}) {}

func (l *NopLogger) Println(v ...interface{}) {}

type LogWriter struct {
	io.Writer
	TimeFormat string
}

func (w LogWriter) Write(b []byte) (n int, err error) {
	return w.Writer.Write(append([]byte(time.Now().Format(w.TimeFormat)), b...))
}
