package demo

import (
	"github.com/kataras/iris/v12"
	"github.com/yuanJewel/go-core/api"
)

const prefixPath = ""

func Party(app iris.Party, topPath string) {
	app.Get(api.Fusion(topPath, prefixPath, "project"), getProjects).Name = "get-projects"
	app.Put(api.Fusion(topPath, prefixPath, "project"), putProject).Name = "put-project"
	app.Post(api.Fusion(topPath, prefixPath, "project"), postProjects).Name = "post-projects"
	app.Delete(api.Fusion(topPath, prefixPath, "project"), deleteProjects).Name = "delete-projects"
}
