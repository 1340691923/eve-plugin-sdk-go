package logger

import "github.com/hashicorp/go-hclog"

type Level int32

const (
	NoLevel	Level	= iota
	Trace
	Debug
	Info
	Warn
	Error
)

type Logger interface {
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Level() Level
}

func New() Logger {
	return NewWithLevel(Debug)
}

func NewWithLevel(level Level) Logger {
	return &hclogWrapper{
		logger: hclog.New(&hclog.LoggerOptions{

			Level:	hclog.Level(level),

			JSONFormat:	true,
		}),
	}
}

type hclogWrapper struct {
	logger hclog.Logger
}

func (l *hclogWrapper) Debug(msg string, args ...interface{}) {
	l.logger.Debug(msg, args...)
}

func (l *hclogWrapper) Info(msg string, args ...interface{}) {
	l.logger.Info(msg, args...)
}

func (l *hclogWrapper) Warn(msg string, args ...interface{}) {
	l.logger.Warn(msg, args...)
}

func (l *hclogWrapper) Error(msg string, args ...interface{}) {
	l.logger.Error(msg, args...)
}

func (l *hclogWrapper) Level() Level {
	if l.logger.IsDebug() {
		return Debug
	}
	if l.logger.IsTrace() {
		return Trace
	}
	if l.logger.IsInfo() {
		return Info
	}
	if l.logger.IsWarn() {
		return Warn
	}
	if l.logger.IsError() {
		return Error
	}
	return NoLevel
}

var DefaultLogger = New()
