package ssql

import (
	"fmt"
)

const (
	LevelError = iota + 1 // 1
	LevelInfo             // 2
	LevelWarn             // 3
	LevelDebug            // 4
)

type Logger interface {
	Print(v ...any)
	Printf(format string, v ...any)
}

type defaultLogger struct{}

func (_ defaultLogger) Print(v ...any) {
	fmt.Print(v...)
}

func (_ defaultLogger) Printf(format string, v ...any) {
	fmt.Printf(format, v...)
}

type log struct {
	logger *Logger
	level  int
}

func (log *log) debug(v ...any) {
	if log.logger == nil || log.level >= 4 {
		return
	}
	(*log.logger).Print(v...)
}

func (log *log) debugf(format string, v ...any) {
	if log.logger == nil || log.level >= 4 {
		return
	}
	(*log.logger).Printf(format, v...)
}

func (log *log) warn(v ...any) {
	if log.logger == nil || log.level >= 3 {
		return
	}
	(*log.logger).Print(v...)
}

func (log *log) warnf(format string, v ...any) {
	if log.logger == nil || log.level >= 3 {
		return
	}
	(*log.logger).Printf(format, v...)
}

func (log *log) info(v ...any) {
	if log.logger == nil || log.level >= 2 {
		return
	}
	(*log.logger).Print(v...)
}

func (log *log) infof(format string, v ...any) {
	if log.logger == nil || log.level >= 2 {
		return
	}
	(*log.logger).Printf(format, v...)
}

func (log *log) error(err error, v ...any) {
	if log.logger == nil || log.level >= 1 {
		return
	}
	v = append(v, fmt.Sprintf("error: %s", err))
	(*log.logger).Print(v...)
}

func (log *log) errorf(err error, format string, v ...any) {
	if log.logger == nil || log.level >= 1 {
		return
	}
	v = append(v, fmt.Sprintf("error: %s", err))
	format = fmt.Sprintf("%s error: %s", format, err)
	(*log.logger).Printf(format, v...)
}
