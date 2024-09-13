package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/SmartLyu/go-core/logger"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/middleware/jwt"
	"net/http"
	"strings"
	"time"
)

// AuthenticateApi 登录认证
// @Summary 登录认证
// @Description 登录认证
// @Param body body AuthenticateConfig true "Account Info"
// @tags login
// @Accept json
// @Produce json
// @Success 200 string api_interface.Response "ok"
// @Failure 401 string api_interface.Response "未授权"
// @Failure 407 string api_interface.Response "权限不足"
// @Failure 501 string api_interface.Response "处理存在异常"
// @Security ApiKeyAuth
// @Router /authenticate [post]
func (Object) AuthenticateApi(ctx iris.Context) {
	response := ResponseInit(ctx)
	url := fmt.Sprintf("%s/authenticate", strings.TrimLeft(AuthenticateUrl, "/"))
	headers := http.Header{"traceid": []string{response.TraceId}}
	body, err := ctx.GetBody()
	if err != nil {
		ReturnErr(GetBodyError, ctx, err, response)
		return
	}
	code, returnBody, err := HttpUtil(http.MethodPost, url, HttpReadTimeout*time.Second, headers, bytes.NewReader(body))
	if err != nil {
		ReturnErr(ConnectAuthenticationError, ctx,
			errors.New(fmt.Sprintf("return code is %d, Error is %v", code, err)), response)
		return
	}
	var returnObject Response
	if err := json.Unmarshal(returnBody, &returnObject); err != nil {
		ReturnErr(JsonUnmarshalError, ctx, err, response)
		return
	}
	ctx.StatusCode(code)
	_ = ctx.JSON(returnObject)
}

func (Object) Authenticate(ctx iris.Context) error {
	url := fmt.Sprintf("%s/api/v1/free/authenticate", strings.TrimLeft(AuthenticateUrl, "/"))
	headers := ctx.Request().Header
	headers.Set("Path", ctx.Path())
	headers.Set("Method", ctx.Method())
	code, returnBody, err := HttpUtil(http.MethodGet, url, HttpReadTimeout*time.Second, headers, nil)
	if code != 200 || err != nil {
		return errors.New(fmt.Sprintf("return code is %d, Error is %v", code, err))
	}
	var (
		returnObject AuthenticationResponse
	)
	if err = json.Unmarshal(returnBody, &returnObject); err != nil {
		return err
	}
	if returnObject.Code != 0 || returnObject.Message != "" {
		return errors.New("permission deny")
	}
	returnData := returnObject.Data

	if !returnData.Approved {
		return errors.New("permission deny")
	}
	for k, v := range returnData.Header {
		ctx.Request().Header.Set(k, v)
	}
	return nil
}

func ParseToken(ctx iris.Context) (string, float64, float64, error) {
	header := strings.TrimPrefix(ctx.GetHeader("Authorization"), "Bearer ")
	if strings.HasPrefix(header, "Basic") {
		return "", 0, 0, errors.New("Basic veritfy cannot get token")
	}
	verifier := jwt.NewVerifier(jwt.HS256, Auth)
	token, err := verifier.VerifyToken([]byte(header))
	if err != nil || token == nil {
		return "", 0, 0, err
	}
	var v interface{}
	if err = json.Unmarshal(token.Payload, &v); err != nil {
		return "", 0, 0, err
	}
	data, ok := v.(map[string]interface{})
	if !ok {
		return "", 0, 0, errors.New("cannot verify token")
	}
	name, ok := data["username"].(string)
	if !ok {
		return "", 0, 0, errors.New("cannot verify token")
	}
	iat, ok := data["iat"].(float64)
	if !ok {
		return "", 0, 0, errors.New("cannot verify token")
	}
	exp, ok := data["exp"].(float64)
	if !ok {
		return "", 0, 0, errors.New("cannot verify token")
	}
	return name, iat, exp, nil
}

func GetUserName(ctx iris.Context) string {
	user, _, _, err := ParseToken(ctx)
	if err != nil {
		logger.Log.Warningln(err)
		return "unknown"
	}
	return user
}

func irisAuthenticate(check func(iris.Context) error) iris.Handler {
	return func(ctx *context.Context) {
		response := ResponseInit(ctx)
		user, _, _, err := ParseToken(ctx)
		if err != nil {
			ReturnErr(ParseTokenEorror, ctx, err, response)
			return
		}
		if user == "admin" {
			ctx.Next()
			return
		}
		path := ctx.Path()
		if strings.HasPrefix(path, "/api/v1/free") {
			ctx.Next()
			return
		}
		if err := check(ctx); err != nil {
			ReturnErr(PermissionDeny, ctx, err, response)
			return
		} else {
			ctx.Next()
		}
	}
}
