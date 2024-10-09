package btlLog

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

type LogLevel int

const (
	ERROR LogLevel = iota
	WARNING
	DEBUG
	INFO
)

func InitBtlLog() error {
	if err := openLogFile(); err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}
	loadDefaultLog()
	return nil
}

type ServicesLogger struct {
	logger      *log.Logger
	errorLogger *log.Logger
	level       LogLevel
}

func NewLogger(logName string, level LogLevel, Writer ...io.Writer) *ServicesLogger {
	var multiWriter io.Writer
	multiWriter = io.MultiWriter(os.Stdout)
	for i := range Writer {
		multiWriter = io.MultiWriter(multiWriter, Writer[i])
	}
	return &ServicesLogger{
		logger:      log.New(multiWriter, "["+logName+"]: ", log.Ldate|log.Ltime),
		errorLogger: log.New(io.MultiWriter(multiWriter, defaultErrorLogFile), "["+logName+"]: ", log.Ldate|log.Ltime),
		level:       level,
	}
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
		ml.errorLogger.Printf(" %s：%d [Error]: %s\n", callerFile, callerLine, msg)
	}
}

var (
	defaultLogFile      *os.File
	defaultErrorLogFile *os.File
)

func openLogFile() error {
	var err error
	dirPath := "./logs"
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		// 如果目录不存在，创建目录
		err := os.MkdirAll(dirPath, 0755)
		if err != nil {
			return err
		}
		fmt.Println("目录已创建:", dirPath)
	}
	filePath := filepath.Join(dirPath, "output.log")
	backupLogFile(filePath)
	defaultLogFile, err = os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	ErrorFilePath := filepath.Join(dirPath, "error.log")
	backupLogFile(ErrorFilePath)
	defaultErrorLogFile, err = os.OpenFile(ErrorFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	return nil
}

func backupLogFile(filePath string) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return
	}
	if fileInfo.Size() > 20*1024*1024 {
		newName := filePath + "." + time.Now().Format("200601021504") + ".bak"
		err := os.Rename(filePath, newName)
		if err != nil {
			fmt.Printf("Backup file failed: %v", err)
		}
	}
}

var (
	CUST                  *ServicesLogger
	FairLaunchDebugLogger *ServicesLogger
	FEE                   *ServicesLogger
	ScheduledTask         *ServicesLogger
	PreSale               *ServicesLogger
)

func loadDefaultLog() {
	Level := INFO
	CUST = NewLogger("CUST", Level, defaultLogFile)
	FairLaunchDebugLogger = NewLogger("FLDL", Level, defaultLogFile)
	FEE = NewLogger("FEE", Level, defaultLogFile)
	ScheduledTask = NewLogger("CRON", Level, defaultLogFile)
	PreSale = NewLogger("PRSL", Level, defaultLogFile)
}
