package main

import (
	"os"

	"github.com/pivotal-golang/lager"
)

//enum used to describe the possible values for logging
const (
	DEBUG = "debug"
	INFO  = "info"
	ERROR = "error"
	FATAL = "fatal"
)

//NewLogger initializes a new logger for usb that will log anything above the "level" specified
func NewLogger(level string) lager.Logger {
	var logger = lager.NewLogger("usb")

	var minLogLevel lager.LogLevel
	switch level {
	case DEBUG:
		minLogLevel = lager.DEBUG
	case INFO:
		minLogLevel = lager.INFO
	case ERROR:
		minLogLevel = lager.ERROR
	case FATAL:
		minLogLevel = lager.FATAL
	default:
		minLogLevel = lager.DEBUG
	}

	logger.RegisterSink(lager.NewWriterSink(os.Stdout, minLogLevel))

	logger.Info("Log level set to:", lager.Data{"level": minLogLevel})

	return logger
}
