package logger

import (
	"fangaoxs.com/go-chat/environment"

	"github.com/sirupsen/logrus"
)

func New(env environment.Env) Logger {
	logger := logrus.New()
	logger.SetLevel(string2level(env.LogLevel))
	logger.SetFormatter(&logrus.TextFormatter{})
	logger.SetReportCaller(false)

	return Logger{
		logger: logger.WithFields(logrus.Fields{
			"app":     env.AppName,
			"version": env.AppVersion,
		}),
	}
}

type Logger struct {
	logger *logrus.Entry
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	l.logger.Debugf(format, args...)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	l.logger.Warnf(format, args...)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.logger.Errorf(format, args...)
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
