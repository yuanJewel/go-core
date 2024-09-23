package api

import (
	"github.com/SmartLyu/go-core/api"
	"github.com/SmartLyu/go-core/db/service"
	"github.com/SmartLyu/go-core/pkg/config"
	"github.com/kataras/iris/v12"
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

	app.Get("/project", getProjects).Name = "get-projects"
	app.Put("/project", putProject).Name = "put-project"
	app.Post("/project", postProjects).Name = "post-projects"
	app.Delete("/project", deleteProjects).Name = "delete-projects"

	app.OnErrorCode(iris.StatusNotFound, notFound)
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
	ctx.ViewData("message", "welcome to SmartLyu go-core")
	if err := ctx.View("index.html"); err != nil {
		return
	}
}

func notFound(ctx iris.Context) {
	if err := ctx.View("404.html"); err != nil {
		return
	}
}
