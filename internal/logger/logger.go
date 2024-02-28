package logger

import log "github.com/sirupsen/logrus"

func New(level string) *log.Logger {
	var logLevel log.Level
	switch level {
	case "info":
		logLevel = log.InfoLevel
	case "debug":
		logLevel = log.DebugLevel
	case "error":
		logLevel = log.ErrorLevel
	case "panic":
		logLevel = log.PanicLevel
	case "warn":
		logLevel = log.WarnLevel
	case "fatal":
		logLevel = log.FatalLevel
	default:
		logLevel = log.InfoLevel
	}

	logger := log.New()
	logger.SetLevel(logLevel)
	return logger
}
