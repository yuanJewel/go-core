package task

import (
	"github.com/kataras/iris/v12"
	"github.com/yuanJewel/go-core/api"
	"github.com/yuanJewel/go-core/db/service"
)

func CreateTaskContext(ctx iris.Context) {
	response := api.ResponseInit(ctx)
	dbInstance := service.Instance.WithContext(ctx)
	body, err := ctx.GetBody()
	if err != nil {
		api.ReturnErr(api.GetBodyError, ctx, err, response)
		return
	}
	var job Job

	err = CreateTask(&job, dbInstance, body, api.GetUserName(ctx))
	if err != nil {
		api.ReturnErr(api.CreateTaskError, ctx, err, response)
		return
	}
	api.ResponseBody(ctx, response, job)
}
