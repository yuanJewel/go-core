package logger

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/yuanJewel/go-core/utils"
	"os"
	"path"
	"strconv"
	"sync"
	"time"
)

var (
	Log = logrus.WithFields(logrus.Fields{"traceId": "-"})
)

func init() {

	Log.Logger.SetReportCaller(true)
	Log.Logger.SetLevel(logrus.InfoLevel)
	Log.Logger.SetFormatter(&logrus.JSONFormatter{
		//DisableColors:   true,
		//FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05,000",
	})
	if !isLogOutFile() {
		Log.Logger.Out = os.Stdout
		return
	}
	_logDir := getLogRootPath()
	if err := Exists(_logDir); err != nil {
		_ = os.MkdirAll(_logDir, os.FileMode(0777))
	}
	_filename := GetLogFilename()
	file, err := os.OpenFile(_filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		Log.Logger.Out = file
	} else {
		Log.Error(err)
		Log.Info("Failed to logger to file, using default stderr")
	}
	go RollLogFile()
}

// RollLogFile 滚动更新日志文件
// 按大小进行滚动，50M 一个文件，默认保留三个文件
// 保留个数可以通过 环境变量 LoggerRetainNumber 覆盖
// 单个日志文件大小 通过 LoggerFileSize 覆盖
func RollLogFile() {
	ticker := time.NewTicker(10 * time.Second)
	retainNumber := getLoggerRetainNumber()
	filename := GetLogFilename()
	defer ticker.Stop()
	for {
		<-ticker.C
		file, err := os.Stat(filename)
		if err != nil {
			continue
		}
		size := utils.BitToMb(file.Size())
		if size > getLogFileSize() {
			rollFile(filename, 0, retainNumber)
		}
	}
}
func rollFile(_filename string, num, retainNumber int) {
	filename := _filename
	if num > 0 {
		filename = fmt.Sprintf("%s.%d", _filename, num)
	}
	if utils.FileExists(filename) {
		if num >= retainNumber {
			_ = os.Remove(filename)
			return
		}
		rollFile(_filename, num+1, retainNumber)
	}
	if num > 0 {
		_ = os.Rename(filename, fmt.Sprintf("%s.%d", _filename, num+1))
	} else {
		var mutex sync.Mutex
		mutex.Lock()
		_ = os.Rename(filename, fmt.Sprintf("%s.%d", _filename, num+1))
		if file, err := os.OpenFile(_filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666); err == nil {
			Log.Logger.Out = file
		}
		mutex.Unlock()
	}
}

func GetLogFilename() string {
	filename := os.Getenv("Appname")
	if filename == "" {
		filename = "golang"
	}
	filename += ".log"
	return path.Join(getLogRootPath(), filename)
}

func getLoggerRetainNumber() int {
	retainNumber := 3
	envNumber := os.Getenv("LoggerRetainNumber")
	if envNumber != "" {
		if n, err := strconv.Atoi(envNumber); err == nil {
			retainNumber = n
		}
	}
	return retainNumber
}

func getLogFileSize() int64 {
	var size int64 = 50
	_size := os.Getenv("LoggerFileSize")
	if _size != "" {
		if n, err := strconv.ParseInt(_size, 10, 64); err != nil {
			size = n
		}
	}
	return size
}

// 配置log的根目录
// 如果未配置根目录则为当前目录
func getLogRootPath() string {
	return path.Join(os.Getenv("LOGGER_ROOT_PATH"), logDir)
}

func isLogOutFile() bool {
	style := os.Getenv("LOGGER_OUT_STYLE")
	if style != "stdout" {
		return true
	}
	return false
}
