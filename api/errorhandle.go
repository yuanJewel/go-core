package api

import (
	"fmt"
	"github.com/SmartLyu/go-core/logger"
	"github.com/kataras/iris/v12"
	"net/http"
	"runtime"
)

const (
	// data error
	PermissionDeny     = 1
	JsonUnmarshalError = 2

	// http error
	ParseTokenEorror = 11
	GetBodyError     = 12
	GetTokenError    = 13
	GetHeaderError   = 14

	// login error
	AuthenticationError = 31
	LdapConnectError    = 32
	LdapSearchError     = 33
	LdapUserAuthError   = 34
	ParseHeaderError    = 35
	GoogleCodeError     = 36

	// task error
	InitWorkerError  = 41
	StartWorkerError = 42
	GetWorkerError   = 43

	// connect other components error
	ConnectAuthenticationError = 101
	ConnectAssetRecordError    = 102
	ConnectTaskWorkerError     = 102

	// DB connnect error
	SelectDbError = 301
	AddDbError    = 302
	UpdateDbError = 303
	DeleteDbError = 304

	// Redis connect error
	ConnectRedisError = 401
	SelectRedisError  = 402

	// reflect error
	ReflectError           = 501
	UnmarshalResponseError = 502
	SpecialReturnError     = 503
)

func ReturnErr(code int, ctx iris.Context, err error, response *Response) {
	var (
		functionName = "unknown_function"
		functionFile = ""
		functionLine = 0
	)
	pc, pcFile, pcLine, ok := runtime.Caller(1)
	if ok {
		functionName = runtime.FuncForPC(pc).Name()
		functionFile = pcFile
		functionLine = pcLine
	}
	errMsg := fmt.Sprintf("%s", err.Error())
	logger.Log.WithField("traceId", response.TraceId).WithField("function", functionName).
		WithField("callerFile", functionFile).WithField("callerLine", functionLine).Warnf(errMsg)
	response.Code = code

	if code == PermissionDeny {
		ctx.StatusCode(http.StatusForbidden)
	} else {
		ctx.StatusCode(http.StatusNotImplemented)
	}
	ResponseBody(ctx, response, errMsg)
}
