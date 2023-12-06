package logger

import (
	"fangaoxs.com/go-chat/environment"

	"github.com/sirupsen/logrus"
)

type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})

	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

func New(env environment.Env) Logger {
	log := logrus.New()
	log.SetLevel(string2level(env.LogLevel))
	log.SetFormatter(&logrus.TextFormatter{})
	log.SetReportCaller(false)

	return &logger{
		log: log.WithFields(logrus.Fields{
			"app":     env.AppName,
			"version": env.AppVersion,
		}),
	}
}

type logger struct {
	log *logrus.Entry
}

func (l *logger) Debug(args ...interface{}) {
	l.log.Debug(args...)
}

func (l *logger) Info(args ...interface{}) {
	l.log.Debug(args...)
}

func (l *logger) Warn(args ...interface{}) {
	l.log.Debug(args...)
}

func (l *logger) Error(args ...interface{}) {
	l.log.Debug(args...)
}

func (l *logger) Debugf(format string, args ...interface{}) {
	l.log.Debugf(format, args...)
}

func (l *logger) Infof(format string, args ...interface{}) {
	l.log.Infof(format, args...)
}

func (l *logger) Warnf(format string, args ...interface{}) {
	l.log.Warnf(format, args...)
}

func (l *logger) Errorf(format string, args ...interface{}) {
	l.log.Errorf(format, args...)
}

func string2level(s string) logrus.Level {
	switch s {
	case "TRACE", "Trace", "trace":
		return logrus.TraceLevel
	case "DEBUG", "Debug", "debug":
		return logrus.DebugLevel
	case "INFO", "Info", "info":
		return logrus.InfoLevel
	case "WARN", "Warn", "warn":
		return logrus.WarnLevel
	case "ERROR", "Error", "error":
		return logrus.ErrorLevel
	case "FATAL", "Fatal", "fatal":
		return logrus.FatalLevel
	case "PANIC", "Panic", "panic":
		return logrus.PanicLevel
	}

	return logrus.InfoLevel
}
