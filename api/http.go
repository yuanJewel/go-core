package api

import (
	"encoding/json"
	"github.com/kataras/iris/v12"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	HttpTimeout     = 10
	HttpReadTimeout = 900
)

func HttpUtil(method, url string, timeout time.Duration, headers http.Header, body io.Reader) (int, []byte, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return 0, nil, err
	}
	req.Header = headers
	client := &http.Client{Timeout: timeout * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	_body, err := ioutil.ReadAll(resp.Body)
	return resp.StatusCode, _body, err
}

func ReverserUtil(ctx iris.Context, method, path string) (int, []byte) {
	commandPath := Reverser.Path(path)
	rec := ctx.Recorder()
	ctx.Exec(method, commandPath)
	code := ctx.GetStatusCode()
	body := rec.Body()
	rec.ResetBody()
	return code, body
}

func UnmarshalResponse(body []byte) (*Response, error) {
	var returnObject Response
	err := json.Unmarshal(body, &returnObject)
	return &returnObject, err
}

func httpMethodUtil(method, url, username, password string, timeout time.Duration, headers map[string]string, body io.Reader) (int, []byte, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return 0, nil, err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	if username != "" && password != "" {
		req.SetBasicAuth(username, password)
	}
	client := &http.Client{Timeout: timeout * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	_body, err := ioutil.ReadAll(resp.Body)
	return resp.StatusCode, _body, err
}
