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
)

func ApiReturnErr(code int, ctx iris.Context, err error, response *Response) {
	var (
		function_name = "unknown_function"
		function_file = ""
		function_line = 0
	)
	pc, _pc_file, _pc_line, ok := runtime.Caller(1)
	if ok {
		function_name = runtime.FuncForPC(pc).Name()
		function_file = _pc_file
		function_line = _pc_line
	}
	errMsg := fmt.Sprintf("%s", err.Error())
	logger.Log.WithField("traceid", response.TraceId).WithField("function", function_name).
		WithField("callerfile", function_file).WithField("callerline", function_line).Warnf(errMsg)
	response.Code = code
	response.Message = errMsg
	response.Data = nil
	if code == PermissionDeny {
		ctx.StatusCode(http.StatusForbidden)
	} else {
		ctx.StatusCode(http.StatusNotImplemented)
	}
	_ = ctx.JSON(response)
}
