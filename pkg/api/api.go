package api

import (
	apiInterface "github.com/SmartLyu/go-core/api"
	"github.com/SmartLyu/go-core/cmdb"
	"github.com/SmartLyu/go-core/pkg/config"
	"github.com/kataras/iris/v12"
)

type Object struct {
	apiInterface.Object
}

func (Object) GetAuth() string {
	return config.GlobalConfig.Auth.Key
}

func (Object) Party(app iris.Party) {
	app.Get("/", index).Name = "index"
}

func (Object) Health() func() map[string]error {
	return func() map[string]error {
		errList := map[string]error{}
		errList["cmdb"] = cmdb.CmdbInstance.HealthCheck()
		return errList
	}
}

func index(ctx iris.Context) {
	ctx.ViewData("message", "welcome to SmartLyu go-core")
	if err := ctx.View("index.html"); err != nil {
		return
	}
}
