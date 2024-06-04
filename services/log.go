package services

import (
	"fmt"
	"log"
	"os"
	"runtime"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARNING
	ERROR
)

type ServicesLogger struct {
	logger *log.Logger
	level  LogLevel
}

func NewLogger(logName string, level LogLevel) *ServicesLogger {
	return &ServicesLogger{
		logger: log.New(os.Stdout, "["+logName+"]: ", log.Ldate|log.Ltime),
		level:  level}
}

func (ml *ServicesLogger) Debug(format string, v ...interface{}) {
	if ml.level >= DEBUG {
		_, callerFile, callerLine, _ := runtime.Caller(1)
		msg := fmt.Sprintf(format, v...)
		ml.logger.Printf(" %s：%d [Debug]: %s\n", callerFile, callerLine, msg)
	}
}

func (ml *ServicesLogger) Info(format string, v ...any) {
	if ml.level >= INFO {
		_, callerFile, callerLine, _ := runtime.Caller(1)
		msg := fmt.Sprintf(format, v...)
		ml.logger.Printf(" %s：%d [Log]: %s\n", callerFile, callerLine, msg)
	}
}

func (ml *ServicesLogger) Warning(format string, v ...interface{}) {
	if ml.level >= WARNING {
		_, callerFile, callerLine, _ := runtime.Caller(1)
		msg := fmt.Sprintf(format, v...)
		ml.logger.Printf(" %s：%d [Warning]: %s\n", callerFile, callerLine, msg)
	}
}

func (ml *ServicesLogger) Error(format string, v ...any) {
	if ml.level >= ERROR {
		_, callerFile, callerLine, _ := runtime.Caller(1)
		msg := fmt.Sprintf(format, v...)
		ml.logger.Printf(" %s：%d [Error]: %s\n", callerFile, callerLine, msg)
	}
}

var (
	CUST                  = NewLogger("CUST", ERROR)
	FairLaunchDebugLogger = NewLogger("FLDL", ERROR)
	FEE                   = NewLogger("FEE", ERROR)
)
