package mysql

import (
	"context"
	"errors"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/sirupsen/logrus"
	"github.com/yuanJewel/go-core/api"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
	"time"
)

const (
	logTraceStr     = "[%.3fms] [rows:%v] %s"
	logTraceWarnStr = "%s [%.3fms] [rows:%v] %s"
	logTraceErrStr  = "%s [%.3fms] [rows:%v] %s"
)

func getTraceId(ctx context.Context) string {
	if c, ok := ctx.(iris.Context); ok {
		return api.GetTraceId(c)
	}
	return "-"
}

type logger struct {
	log *logrus.Entry
	gormlogger.Config
}

// LogMode log mode
func (l *logger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	newLogger := *l
	newLogger.LogLevel = level
	return &newLogger
}

// Info print info
func (l *logger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlogger.Info {
		l.log.WithField("traceId", getTraceId(ctx)).Infof(msg, data...)
	}
}

// Warn print warn messages
func (l *logger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlogger.Warn {
		l.log.WithField("traceId", getTraceId(ctx)).Warnf(msg, data...)
	}
}

// Error print error messages
func (l *logger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlogger.Error {
		l.log.WithField("traceId", getTraceId(ctx)).Errorf(msg, data...)
	}
}

// Trace print sql message
func (l *logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= gormlogger.Silent {
		return
	}

	log := l.log.WithField("traceId", getTraceId(ctx)).WithField("callerFile", utils.FileWithLineNum())
	elapsed := time.Since(begin)
	switch {
	case err != nil && l.LogLevel >= gormlogger.Error && (!errors.Is(err, gormlogger.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError):
		sql, rows := fc()
		if rows == -1 {
			log.Errorf(logTraceErrStr, err, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			log.Errorf(logTraceErrStr, err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= gormlogger.Warn:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL OVER (%v)", l.SlowThreshold)
		if rows == -1 {
			log.Warnf(logTraceWarnStr, slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			log.Warnf(logTraceWarnStr, slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case l.LogLevel == gormlogger.Info:
		sql, rows := fc()
		if rows == -1 {
			log.Debugf(logTraceStr, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			log.Debugf(logTraceStr, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	}
}

type nilLogger struct {
	log *logrus.Entry
	gormlogger.Config
}

func (l *nilLogger) LogMode(gormlogger.LogLevel) gormlogger.Interface {
	return l
}
func (l *nilLogger) Info(context.Context, string, ...interface{})  {}
func (l *nilLogger) Warn(context.Context, string, ...interface{})  {}
func (l *nilLogger) Error(context.Context, string, ...interface{}) {}
func (l *nilLogger) Trace(context.Context, time.Time, func() (sql string, rowsAffected int64), error) {
}
