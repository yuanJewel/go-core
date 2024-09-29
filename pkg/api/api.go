package api

import (
	"github.com/kataras/iris/v12"
	"github.com/yuanJewel/go-core/api"
	"github.com/yuanJewel/go-core/db/service"
	"github.com/yuanJewel/go-core/pkg/api/demo"
	"github.com/yuanJewel/go-core/pkg/config"
)

type Object struct {
	api.Object
}

func (Object) GetAuth() string {
	return config.GlobalConfig.Auth.Key
}

func (Object) Party(app iris.Party) {
	app.Get("/", index).Name = "index"
	app.Get("/free/refresh", refresh).Name = "refresh"
	app.OnErrorCode(iris.StatusNotFound, notFound)

	demo.Party(app, "/")
}

func (Object) AuthenticateApi(ctx iris.Context) {
	authenticate(ctx)
}

func (Object) Health() func() map[string]error {
	return func() map[string]error {
		errList := map[string]error{}
		errList["cmdb"] = service.Instance.HealthCheck()
		return errList
	}
}

func index(ctx iris.Context) {
	ctx.ViewData("message", "welcome to yuanJewel go-core")
	if err := ctx.View("index.html"); err != nil {
		return
	}
}

func notFound(ctx iris.Context) {
	if err := ctx.View("404.html"); err != nil {
		return
	}
}
