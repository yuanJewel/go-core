package api

import (
	"github.com/SmartLyu/go-core/api"
	"github.com/SmartLyu/go-core/db/service"
	"github.com/SmartLyu/go-core/pkg/db"
	"github.com/google/uuid"
	"github.com/kataras/iris/v12"
)

// @Summary 获取项目信息
// @Description 获取项目信息
// @Param ids header string true "project"
// @Param body body db.ProjectInfo false "Info"
// @tags project
// @Accept json
// @Produce json
// @Success 200 {object} api.Response "ok"
// @Failure 401 string string "未授权"
// @Failure 403 {object} api.Response "权限不足"
// @Failure 501 {object} api.Response "处理存在异常"
// @Security ApiKeyAuth
// @Router /api/v1/project [get]
func getProjects(ctx iris.Context) {
	service.GetDbInfoByIds(ctx, db.Project{}, &[]db.Project{})
}

// @Summary 新增项目信息
// @Description 新增项目信息
// @Param body body []db.ProjectInfo true "Info"
// @tags project
// @Accept json
// @Produce json
// @Success 200 {object} api.Response "ok"
// @Failure 401 string string "未授权"
// @Failure 403 {object} api.Response "权限不足"
// @Failure 501 {object} api.Response "处理存在异常"
// @Security ApiKeyAuth
// @Router /api/v1/project [post]
func postProjects(ctx iris.Context) {
	service.PostDbInfo(ctx, &[]db.Project{}, func(m *map[string]interface{}) error {
		if err := api.NormalSpecialTask(m); err != nil {
			return err
		}
		(*m)["id"] = uuid.New().String()
		return nil
	})
}

// @Summary 修改项目信息
// @Description 修改项目信息
// @Param id header string true "project"
// @Param body body db.ProjectInfo true "Info"
// @tags project
// @Accept json
// @Produce json
// @Success 200 {object} api.Response "ok"
// @Failure 401 string string "未授权"
// @Failure 403 {object} api.Response "权限不足"
// @Failure 501 {object} api.Response "处理存在异常"
// @Security ApiKeyAuth
// @Router /api/v1/project [put]
func putProject(ctx iris.Context) {
	service.PutDbInfoById(ctx, "get-projects", &db.Project{}, api.NormalSpecialTask)
}

// @Summary 删除项目信息
// @Description 删除项目信息
// @Param ProjectIds header string true "project"
// @tags project
// @Accept json
// @Produce json
// @Success 200 {object} api.Response "ok"
// @Failure 401 string string "未授权"
// @Failure 403 {object} api.Response "权限不足"
// @Failure 501 {object} api.Response "处理存在异常"
// @Security ApiKeyAuth
// @Router /api/v1/project [delete]
func deleteProjects(ctx iris.Context) {
	service.DeleteDb(ctx, "get-projects", &db.Project{})
}
