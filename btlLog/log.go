package btlLog

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"
	"trade/utils"
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

func NewLogger(logName string, level LogLevel, hasStdout bool, Writer ...io.Writer) *ServicesLogger {
	var multiWriter io.Writer

	if hasStdout {
		multiWriter = io.MultiWriter(os.Stdout)
	}
	for i := range Writer {
		if multiWriter == nil {
			multiWriter = io.MultiWriter(Writer[i])
			continue
		}
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
	presaleLogFile      *os.File
	mintNftFile         *os.File
	userDataLogFile     *os.File
	userStatsLogFile    *os.File
	cpAmmLogFile        *os.File
	dateIpLoginLogFile  *os.File
	pushQueueLogFile    *os.File
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
	presaleLogFile, err = utils.GetLogFile("./logs/trade.presale.log")
	if err != nil {
		return err
	}
	mintNftFile, err = utils.GetLogFile("./logs/trade.mint_nft.log")
	if err != nil {
		return err
	}
	userDataLogFile, err = utils.GetLogFile("./logs/trade.userdata.log")
	if err != nil {
		return err
	}
	userStatsLogFile, err = utils.GetLogFile("./logs/trade.user_stats.log")
	if err != nil {
		return err
	}
	cpAmmLogFile, err = utils.GetLogFile("./logs/trade.cp_amm.log")
	if err != nil {
		return err
	}
	dateIpLoginLogFile, err = utils.GetLogFile("./logs/trade.date_ip_login.log")
	if err != nil {
		return err
	}
	pushQueueLogFile, err = utils.GetLogFile("./logs/trade.push_queue.log")
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
	MintNft               *ServicesLogger
	UserData              *ServicesLogger
	UserStats             *ServicesLogger
	CPAmm                 *ServicesLogger
	DateIpLogin           *ServicesLogger
	PushQueue             *ServicesLogger
)

func loadDefaultLog() {
	Level := INFO
	CUST = NewLogger("CUST", Level, true, defaultLogFile)
	FairLaunchDebugLogger = NewLogger("FLDL", Level, true, defaultLogFile)
	FEE = NewLogger("FEE", Level, true, defaultLogFile)
	ScheduledTask = NewLogger("CRON", Level, true, defaultLogFile)
	PreSale = NewLogger("PRSL", Level, true, defaultLogFile, presaleLogFile)
	MintNft = NewLogger("MINT", Level, false, mintNftFile)
	UserData = NewLogger("URDT", Level, true, defaultLogFile, userDataLogFile)
	UserStats = NewLogger("USTS", Level, true, defaultLogFile, userStatsLogFile)
	CPAmm = NewLogger("CPAM", Level, true, defaultLogFile, cpAmmLogFile)
	DateIpLogin = NewLogger("DILR", Level, true, defaultLogFile, dateIpLoginLogFile)
	PushQueue = NewLogger("PUSH", Level, true, defaultLogFile, pushQueueLogFile)
}
