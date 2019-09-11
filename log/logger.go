package log

import "log"

// Logger defines the logger interface for tango use
type Logger interface {
	Debugf(format string, v ...interface{})
	Debug(v ...interface{})
	Infof(format string, v ...interface{})
	Info(v ...interface{})
	Warnf(format string, v ...interface{})
	Warn(v ...interface{})
	Errorf(format string, v ...interface{})
	Error(v ...interface{})
}

func Infof(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func Info(v ...interface{}) {
	log.Print(v...)
}

func Errorf(format string, v ...interface{}) {
	log.Printf(format, v...)
}
func Error(v ...interface{}) {
	log.Print(v...)
}
