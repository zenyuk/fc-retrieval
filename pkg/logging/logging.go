package logging

import (
	"fmt"
	"log"
	"os"
	"sync"
)

const (
	logLevelStrTest = "TEST"
	// LogLevelStrNone log level should be used when no logging should be emitted
	LogLevelStrNone = "NONE"
	// LogLevelStrError log level should be used when only Error level logging should be emitted
	LogLevelStrError = "ERROR"
	// LogLevelStrWarn log level should be used when only Error and Warn level logging should be emitted
	LogLevelStrWarn = "WARN"
	// LogLevelStrInfo log level should be used when only Error, Warn, and Info level logging should be emitted
	LogLevelStrInfo = "INFO"
	// LogLevelStrTrace log level should be used when all logging should be emitted
	LogLevelStrTrace = "TRACE"
	
	// LogLevelNone log level should be used when no logging should be emitted
	LogLevelNone = 0
	// LogLevelError log level should be used when only Error level logging should be emitted
	LogLevelError = 1
	// LogLevelWarn log level should be used when only Error and Warn level logging should be emitted
	LogLevelWarn = 2
	// LogLevelInfo log level should be used when only Error, Warn, and Info level logging should be emitted
	LogLevelInfo = 3
	// LogLevelTrace log level should be used when all logging should be emitted
	LogLevelTrace = 4

	logTargetStdOut = "STDOUT"
	logTargetStdOutInt = 0
	logTargetFile = "FILE"
	logTargetFileInt = 1

	logFileName = "/etc/gateway/gateway.log"

)


var logLevel = int(LogLevelTrace)
var logTarget = int(logTargetStdOutInt)
var mutex sync.RWMutex

// Init inialises the logging system
func Init() {
	mutex = sync.RWMutex{}
}

// SetLogLevel allows the log level to be specified.
func SetLogLevel(level string) {
	switch (level) {
	case LogLevelStrError:
		logLevel = LogLevelError
	case LogLevelStrWarn:
		logLevel = LogLevelWarn
	case LogLevelStrInfo:
		logLevel = LogLevelInfo
	case LogLevelStrTrace:
		logLevel = LogLevelTrace
	default:
		panic("Unknown Log Level: " + level)
	}
}

// SetLogTarget allows the destination of logs to be specified.
func SetLogTarget(target string) {
	switch (target) {
	case logTargetStdOut:
		logTarget = logTargetStdOutInt
	case logTargetFile:
		logTarget = logTargetFileInt
	default:
		panic("Unknown Log Target: " + target)
	}
}

// ErrorEnabled returns true if ERROR log level is enabled.
func ErrorEnabled() bool {
	return LogLevelError <= logLevel
}

// WarnEnabled returns true if WARN log level is enabled.
func WarnEnabled() bool {
	return LogLevelWarn <= logLevel
}

// InfoEnabled returns true if INFO log level is enabled.
func InfoEnabled() bool {
	return LogLevelInfo <= logLevel
}

// TraceEnabled returns true if TRACE log level is enabled.
func TraceEnabled() bool {
	return LogLevelTrace <= logLevel
}

// Test prints out msg to the log target independant of the log level. 
// This function should only be called from test code.
// The idea is that using the logging framework will put messages from test
// code in the context of the other messages that are emitted.
// msg is interpreted as a format string and args as parameters to the format
// string is there are any args.
func Test(msg string, args ...interface{}) {
	printf(logLevelStrTest, msg, args...)
}


// Error prints out msg to the log target if the log level is ERROR or lower. 
// msg is interpreted as a format string and args as parameters to the format
// string is there are any args.
func Error(msg string, args ...interface{}) {
	if (LogLevelError <= logLevel) {
		printf(LogLevelStrError, msg, args...)
	}
}

// Error1 prints out msg to the log target if the log level is ERROR or lower. 
func Error1(err error) {
	if (LogLevelError <= logLevel) {
		printf(LogLevelStrError, err.Error())
	}
}

// ErrorAndPanic logs an error and then calls panic with the same message.
func ErrorAndPanic(msg string, args ...interface{}) {
	if (LogLevelError <= logLevel) {
		printf(LogLevelStrError, msg, args...)
	}
	panic(fmt.Sprintf(msg, args...))
}

// Warn prints out msg to the log target if the log level is WARN or lower. 
// msg is interpreted as a format string and args as parameters to the format
// string is there are any args.
func Warn(msg string, args ...interface{}) {
	if (LogLevelWarn <= logLevel) {
		printf(LogLevelStrWarn, msg, args...)
	}
}

// Info prints out msg to the log target if the log level is INFO or lower. 
// msg is interpreted as a format string and args as parameters to the format
// string is there are any args.
func Info(msg string, args ...interface{}) {
	if (LogLevelInfo <= logLevel) {
		printf(LogLevelStrInfo, msg, args...)
	}
}

// Trace prints out msg to the log target if the log level is INFO or lower. 
// msg is interpreted as a format string and args as parameters to the format
// string is there are any args.
func Trace(msg string, args ...interface{}) {
	if (LogLevelTrace <= logLevel) {
		printf(LogLevelStrTrace, msg, args...)
	}
}

func logLevelEnabled(level int) bool {
	return level <= logLevel
}


func printf(level string, msg string, args ...interface{}) {
	var s string
	if (len(args) > 0) {
		s = fmt.Sprintf(msg, args...)
	} else {
		s = msg
	}
	s = level + ": " + s

	// TODO will need to do better than this in the server so we don't block all of the time!
	// Block here to ensure only one thread at a time writes to the output device.
	mutex.Lock()
	defer mutex.Unlock()

	switch (logTarget) {
	case logTargetStdOutInt:
		// TODO will need to do better than this in the server so we don't block all of the time!
		log.Println(s)
	case logTargetFileInt:
		f, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		if _, err := f.WriteString(s + "\n"); err != nil {
			panic(err)
		}

	default:
		panic("Unknown log target")
	}
}