package logger

import (
	"fmt"
	"github.com/SmartLyu/go-core/utils"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/robfig/cron"
	"github.com/ryanuber/columnize"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

const deleteFileOnExit = false
const logDir = "logs"

var logFile *os.File

func todayFilename() string {
	today := time.Now().Format("2006-01-02")
	return fmt.Sprintf("%s/access.%s.log", getLogRootPath(), today)
}

func newLogFile() *os.File {
	_logdir := getLogRootPath()
	if err := Exists(_logdir); err != nil {
		if err := os.MkdirAll(_logdir, os.FileMode(0777)); err != nil {
			fmt.Println(err)
		}
	}
	filename := todayFilename()
	//打开一个输出文件，如果重新启动服务器，它将追加到今天的文件中
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	return f
}

var excludeExtensions = [...]string{
	".js",
	".css",
	".jpg",
	".png",
	".ico",
	".svg",
	".html",
	".json",
}

func Columnize(nowFormatted, traceId, username string, latency time.Duration, status int, ip, method, path string, _byte int, message interface{}, headerMessage interface{}) string {
	// 时间 traceId ip method path response time byte useragent
	line := fmt.Sprintf("%s %s %s %s %s %s %v %v %d %s", nowFormatted, traceId, username, ip, method, path, status, int64(latency/time.Millisecond), _byte, headerMessage)
	outputC := []string{
		line,
	}
	output := columnize.SimpleFormat(outputC) + "\n"
	return output
}

func NewRequestLogger(dot func(...interface{})) (h iris.Handler, close func() error) {
	close = func() error { return nil }
	c := logger.Config{
		Status:            true,
		IP:                true,
		Method:            true,
		Path:              true,
		MessageHeaderKeys: []string{"User-Agent", "traceId"},
	}
	logFile = newLogFile()
	close = func() error {
		err := logFile.Close()
		if deleteFileOnExit {
			err = os.Remove(logFile.Name())
		}
		return err
	}
	c.LogFuncCtx = func(ctx iris.Context, latency time.Duration) {
		traceId := ctx.Request().Header.Get("traceId")
		if traceId == "" {
			traceId = GetTraceId(ctx)
		}
		username := ctx.Request().Header.Get("username")
		if username == "" {
			username = "-"
		}
		username, err := url.QueryUnescape(username)
		if err != nil {
			username = "-"
		}
		output := Columnize(fmt.Sprintf("[%s]", time.Now().Format("02/Jan/2006:15:04:05 -0700")), traceId, username, latency, ctx.ResponseWriter().StatusCode(), ctx.RemoteAddr(), ctx.Method(), ctx.Path(), ctx.ResponseWriter().Written(), "", ctx.Request().Header.Get("User-Agent"))
		_, _ = logFile.Write([]byte(output))
		if dot != nil {
			go dot(traceId, username, ctx)
		}
	}
	//我们不想使用记录器，一些静态请求等
	c.AddSkipper(func(ctx iris.Context) bool {
		path := ctx.Path()
		for _, ext := range excludeExtensions {
			if strings.HasSuffix(path, ext) {
				return true
			}
		}
		return false
	})
	h = logger.New(c)
	go rollingAccessLog()
	return
}

func Exists(path string) error {
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return err
	}
	return nil
}

func GetTraceId(ctx iris.Context) string {
	req := ctx.Request()
	headers := req.Header
	traceId := headers.Get("traceId")
	if traceId == "" {
		traceId = utils.CreateUUID()
		headers.Set("traceId", traceId)
	}
	return traceId
}

func rollingAccessLog() {
	c := cron.New()
	if err := c.AddFunc("0 0 0 * * *", rollAccess); err != nil {
		fmt.Println("定期清理任务启动失败，", err)
		return
	}
	c.Start()
}

func rollAccess() {
	accessLogs, err := getAccessLogs()
	if err != nil {
		fmt.Println("滚动access日志文件出错:", err)
		return
	}
	var mutex sync.Mutex
	mutex.Lock()
	logFile = newLogFile()
	mutex.Unlock()
	for _, file := range accessLogs {
		fileSplit := strings.Split(file, ".")
		if len(fileSplit) > 2 {
			if timeDate, err := utils.TimeParseYYYYMMDD(fileSplit[1]); err == nil {
				if time.Now().Sub(timeDate).Hours() >= float64(getLoggerRetainNumber()*24) {
					_ = os.Remove(file)
				}
			}
		}
	}
}

func getAccessLogs() ([]string, error) {
	files := []string{}
	files, err := utils.GetAllFile(logDir, files)
	if err != nil {
		return nil, err
	}
	result := []string{}
	for _, file := range files {
		if strings.Contains(file, "access_log") {
			result = append(result, file)
		}
	}
	return result, nil
}
