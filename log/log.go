package log

import (
	"log"
)

type LogLevelT int

var gLogger *Logger

const (
	Debug LogLevelT = iota + 1
	Info
	Warning
	Error
	Fatal
)

type Logger struct {
	logLevel LogLevelT
}

func GetLogLevelStr(logLevel LogLevelT) string {
	switch logLevel {
	case Debug:
		return "DEBUG"
	case Info:
		return "INFO"
	case Warning:
		return "WARN"
	case Error:
		return "ERROR"
	case Fatal:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

func NewLogger(logLevel LogLevelT) {
	logger := &Logger{logLevel: logLevel}
	if gLogger == nil {
		gLogger = logger
	} else {
		log.Fatalf("trying to initialize logger twice")
	}
}

func Logf(logLevel LogLevelT, format string, args ...interface{}) {

	if gLogger == nil {
		log.Fatalf("trying to use uninitialized logger utilities")
	}

	if logLevel < gLogger.logLevel {
		return
	}

	levelPrefix := GetLogLevelStr(logLevel)

	formatStr := "[" + levelPrefix + "] " + format
	if logLevel == Fatal {
		log.Fatalf(formatStr, args...)
	} else {
		log.Printf(formatStr, args...)
	}
}

func Debugf(format string, args ...interface{}) {
	Logf(Debug, format, args...)
}

func Infof(format string, args ...interface{}) {
	Logf(Info, format, args...)
}

func Warningf(format string, args ...interface{}) {
	Logf(Warning, format, args...)
}

func Errorf(format string, args ...interface{}) {
	Logf(Error, format, args...)
}

func Fatalf(format string, args ...interface{}) {
	Logf(Fatal, format, args...)
}
