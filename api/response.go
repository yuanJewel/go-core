package api

import (
	"github.com/google/uuid"
	"github.com/kataras/iris/v12"
)

// Response api 返回的统一数据结构
// traceId 首先从Header中获取，如果有继续往下传递，否则新生成一个traceId
// Code 用于标识返回状态, Code的定义应该全局统一，看Code就知道错误类型，0 为正常
// Message 用于当api调用出错时返回信息
// Data 是正常返回时，返回的数据
type Response struct {
	TraceId string      `json:"traceId"`
	Code    int         `json:"code,omitempty"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// AuthenticationData 鉴权组件返回的Data的内容格式
type AuthenticationData struct {
	Header   map[string]string `json:"header,omitempty"`
	Approved bool              `json:"approved"`
}

type AuthenticationResponse struct {
	Response
	Data AuthenticationData `json:"data,omitempty"`
}

func ResponseInit(ctx iris.Context) (response *Response) {
	req := ctx.Request()
	headers := req.Header
	traceId := headers.Get("traceId")
	if traceId == "" {
		traceId = uuid.New().String()
		headers.Set("traceId", traceId)
	}
	response = &Response{TraceId: traceId}
	return
}

func GetTraceId(ctx iris.Context) string {
	traceId := ctx.Request().Header.Get("traceId")
	if traceId == "" {
		traceId = "-"
	}
	return traceId
}
