package api

import (
	"github.com/kataras/iris/v12"
	"github.com/sirupsen/logrus"
	"github.com/yuanJewel/go-core/logger"
)

func Logger(ctx iris.Context) *logrus.Entry {
	return logger.Log.WithFields(logrus.Fields{"traceId": GetTraceId(ctx)})
}
