package api

import (
	apiInterface "github.com/SmartLyu/go-core/api"
	"github.com/SmartLyu/go-core/cmdb"
	"github.com/SmartLyu/go-core/pkg/db"
	"github.com/google/uuid"
	"github.com/kataras/iris/v12"
)

// @Summary 获取项目信息
// @Description 获取项目信息
// @Param ids header string true "project"
// @tags project
// @Accept json
// @Produce json
// @Success 200 {object} apiInterface.Response "ok"
// @Failure 401 string string "未授权"
// @Failure 403 {object} apiInterface.Response "权限不足"
// @Failure 501 {object} apiInterface.Response "处理存在异常"
// @Security ApiKeyAuth
// @Router /api/v1/project [get]
func getProjects(ctx iris.Context) {
	cmdb.GetDbByIdsInfo(ctx, db.Project{}, &[]db.Project{})
}

// @Summary 新增项目信息
// @Description 新增项目信息
// @Param body body []db.ProjectInfo true "Info"
// @tags project
// @Accept json
// @Produce json
// @Success 200 {object} apiInterface.Response "ok"
// @Failure 401 string string "未授权"
// @Failure 403 {object} apiInterface.Response "权限不足"
// @Failure 501 {object} apiInterface.Response "处理存在异常"
// @Security ApiKeyAuth
// @Router /api/v1/project [post]
func postProject(ctx iris.Context) {
	cmdb.PostDbInfo(ctx, &[]db.Project{}, func(m *map[string]interface{}) {
		(*m)["id"] = uuid.New().String()
	})
}

// @Summary 修改项目信息
// @Description 修改项目信息
// @Param id header string true "project"
// @Param body body db.ProjectInfo true "Info"
// @tags project
// @Accept json
// @Produce json
// @Success 200 {object} apiInterface.Response "ok"
// @Failure 401 string string "未授权"
// @Failure 403 {object} apiInterface.Response "权限不足"
// @Failure 501 {object} apiInterface.Response "处理存在异常"
// @Security ApiKeyAuth
// @Router /api/v1/project [put]
func putProjects(ctx iris.Context) {
	cmdb.PutDbByIdInfo(ctx, &db.Project{}, apiInterface.NonSpecialTask)
}

// @Summary 删除项目信息
// @Description 删除项目信息
// @Param ProjectIds header string true "project"
// @tags project
// @Accept json
// @Produce json
// @Success 200 {object} apiInterface.Response "ok"
// @Failure 401 string string "未授权"
// @Failure 403 {object} apiInterface.Response "权限不足"
// @Failure 501 {object} apiInterface.Response "处理存在异常"
// @Security ApiKeyAuth
// @Router /api/v1/project [delete]
func deleteProjects(ctx iris.Context) {
	cmdb.DeleteDbByIdInfo(ctx, "get-projects", &db.Project{})
}
