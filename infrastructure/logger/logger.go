package logger

import (
	"log"
	"os"
)

type Logger interface {
	Debug(format string, v ...interface{})
	Info(format string, v ...interface{})
	Warn(format string, v ...interface{})
	Error(format string, v ...interface{})
	Fatal(format string, v ...interface{})
}

type StandardLogger struct {
	debugLogger *log.Logger
	infoLogger  *log.Logger
	warnLogger  *log.Logger
	errorLogger *log.Logger
	fatalLogger *log.Logger
}

func NewStandardLogger() *StandardLogger {
	return &StandardLogger{
		debugLogger: log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile),
		infoLogger:  log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		warnLogger:  log.New(os.Stdout, "WARN: ", log.Ldate|log.Ltime|log.Lshortfile),
		errorLogger: log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
		fatalLogger: log.New(os.Stderr, "FATAL: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

func (l *StandardLogger) Debug(format string, v ...interface{}) {
	l.debugLogger.Printf(format, v...)
}

func (l *StandardLogger) Info(format string, v ...interface{}) {
	l.infoLogger.Printf(format, v...)
}

func (l *StandardLogger) Warn(format string, v ...interface{}) {
	l.warnLogger.Printf(format, v...)
}

func (l *StandardLogger) Error(format string, v ...interface{}) {
	l.errorLogger.Printf(format, v...)
}

func (l *StandardLogger) Fatal(format string, v ...interface{}) {
	l.fatalLogger.Fatalf(format, v...)
}
