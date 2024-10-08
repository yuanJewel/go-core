package api

import (
	"bytes"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/jwt"
	"github.com/yuanJewel/go-core/api"
	"github.com/yuanJewel/go-core/db/object"
	"github.com/yuanJewel/go-core/db/service"
	"github.com/yuanJewel/go-core/logger"
	"github.com/yuanJewel/go-core/pkg/config"
	"github.com/yuanJewel/go-core/utils"
	"strings"
	"time"
)

// @Summary 登录认证
// @Description 登录认证
// @Param body body AuthenticateConfig true "Account Info"
// @tags login
// @Accept json
// @Produce json
// @Success 200 {object} api.Response "ok"
// @Failure 401 string string "未授权"
// @Failure 403 {object} api.Response "权限不足"
// @Failure 501 {object} api.Response "处理存在异常"
// @Security ApiKeyAuth
// @Router /authenticate [post]
func authenticate(ctx iris.Context) {
	response := api.ResponseInit(ctx)
	body, err := ctx.GetBody()
	if err != nil {
		api.ReturnErr(api.GetBodyError, ctx, err, response)
		return
	}

	var (
		auth_object AuthenticateConfig
		_user       object.User
	)
	err = json.Unmarshal(body, &auth_object)
	if err != nil {
		api.ReturnErr(api.JsonUnmarshalError, ctx, err, response)
		return
	}

	_user_exist, err := service.Instance.GetItem(object.User{Name: auth_object.Username}, &_user)
	if err != nil {
		api.ReturnErr(api.SelectDbError, ctx, err, response)
		return
	}
	encryptPassword := config.GlobalConfig.Auth.CryptoPrefix + utils.AesEncrypt(auth_object.Password,
		config.GlobalConfig.Auth.CryptoKey)
	if _user.Passwd != encryptPassword {
		api.ReturnErr(api.AuthenticationError, ctx,
			errors.New("account password entered incorrectly"), response)
		return
	}
	if !verifyCode(_user.GoogleSecret, auth_object.GoogleCode) {
		api.ReturnErr(api.GoogleCodeError, ctx,
			errors.New("incorrect google code"), response)
		return
	}

	token, err := generateToken(auth_object.Username)
	if err != nil {
		api.ReturnErr(api.GetTokenError, ctx, err, response)
	}
	returnData := map[string]interface{}{
		"header":   "Bearer ",
		"token":    token,
		"userName": auth_object.Username,
		"deadline": time.Now().Add(time.Duration(config.GlobalConfig.Auth.Timeout) * time.Minute).Unix(),
		"refresh":  time.Now().Add(time.Duration(config.GlobalConfig.Auth.Refresh) * time.Minute).Unix(),
	}
	if _user_exist {
		if _, err := service.Instance.UpdateItem(_user, &object.User{LastLoginTime: time.Now()}, 1); err != nil {
			api.ReturnErr(api.UpdateDbError, ctx, err, response)
			return
		}
	} else {
		returnData["secret"] = _user.GoogleSecret
	}
	api.ResponseBody(ctx, response, returnData)
	logger.Log.Infof("user %s login successfully", auth_object.Username)
}

func generateToken(username string) (string, error) {
	token, err := jwt.NewSigner(jwt.HS256, config.GlobalConfig.Auth.Key,
		time.Duration(config.GlobalConfig.Auth.Timeout)*time.Minute).Sign(api.UserClaims{Username: username})
	return string(token), err
}

func getGoogleSecret() string {
	const dictionary = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	var _bytes = make([]byte, 16)
	_, _ = rand.Read(_bytes)
	for k, v := range _bytes {
		_bytes[k] = dictionary[v%byte(len(dictionary))]
	}
	return string(_bytes)
}

func verifyCode(secret string, code string) bool {
	// now time
	_time := time.Now().Unix() / 30
	if getGoogleToken(secret, _time) == code {
		return true
	}

	// before 30 second
	if getGoogleToken(secret, _time-1) == code {
		return true
	}

	// after 30 second
	if getGoogleToken(secret, _time+1) == code {
		return true
	}

	return false
}

func getGoogleToken(secret string, interval int64) string {
	//Converts secret to base32 Encoding. Base32 encoding desires a 32-character
	//subset of the twenty-six letters A–Z and ten digits 0–9
	key, err := base32.StdEncoding.DecodeString(strings.ToUpper(secret))
	if err != nil {
		logger.Log.Errorln(err)
	}
	bs := make([]byte, 8)
	binary.BigEndian.PutUint64(bs, uint64(interval))

	//Signing the value using HMAC-SHA1 Algorithm
	hash := hmac.New(sha1.New, key)
	hash.Write(bs)
	h := hash.Sum(nil)

	// We're going to use a subset of the generated hash.
	// Using the last nibble (half-byte) to choose the index to start from.
	// This number is always appropriate as it's maximum decimal 15, the hash will
	// have the maximum index 19 (20 bytes of SHA1) and we need 4 bytes.
	o := (h[19] & 15)

	var header uint32
	//Get 32 bit chunk from hash starting at the o
	r := bytes.NewReader(h[o : o+4])
	err = binary.Read(r, binary.BigEndian, &header)
	if err != nil {
		logger.Log.Errorln(err)
	}

	//Ignore most significant bits as per RFC 4226.
	//Takes division from one million to generate a remainder less than < 7 digits
	h12 := (int(header) & 0x7fffffff) % 1000000
	otp := fmt.Sprintf("%06d", h12)

	return otp
}

// @Summary 刷新token
// @Description 刷新token
// @tags login
// @Accept json
// @Produce json
// @Success 200 {object} api.Response "ok"
// @Failure 401 string string "未授权"
// @Failure 403 {object} api.Response "权限不足"
// @Failure 501 {object} api.Response "处理存在异常"
// @Security ApiKeyAuth
// @Router /api/v1/free/refresh [Get]
func refresh(ctx iris.Context) {
	response := api.ResponseInit(ctx)
	user, iat, exp, err := api.ParseToken(ctx)
	if err != nil {
		api.ReturnErr(api.ParseTokenEorror, ctx, err, response)
		return
	}
	refreshTime := time.Unix(int64(iat), 0).Add(time.Duration(config.GlobalConfig.Auth.Refresh) * time.Minute).Unix()
	if refreshTime > time.Now().Unix() {
		response.Data = map[string]interface{}{
			"header":   "Bearer ",
			"token":    strings.TrimPrefix(ctx.GetHeader("Authorization"), "Bearer "),
			"userName": user,
			"deadline": int64(exp),
			"refresh":  refreshTime,
		}
	} else {
		token, err := generateToken(user)
		if err != nil {
			api.ReturnErr(api.GetTokenError, ctx, err, response)
		}
		response.Data = map[string]interface{}{
			"header":   "Bearer ",
			"token":    token,
			"userName": user,
			"deadline": time.Now().Add(time.Duration(config.GlobalConfig.Auth.Timeout) * time.Minute).Unix(),
			"refresh":  time.Now().Add(time.Duration(config.GlobalConfig.Auth.Refresh) * time.Minute).Unix(),
		}
	}
	api.ResponseBody(ctx, response, nil)
}
