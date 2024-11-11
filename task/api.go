package task

import (
	"encoding/json"
	"fmt"
	"github.com/RichardKnop/machinery/v1/tasks"
	"github.com/google/uuid"
	"github.com/kataras/iris/v12"
	"github.com/yuanJewel/go-core/api"
	"github.com/yuanJewel/go-core/db/service"
	"time"
)

func CreateTaskContext(ctx iris.Context) {
	response := api.ResponseInit(ctx)
	instance := service.Instance.WithContext(ctx)
	body, err := ctx.GetBody()
	if err != nil {
		api.ReturnErr(api.GetBodyError, ctx, err, response)
		return
	}
	// 解析 JSON 数据
	var obj interface{}
	err = json.Unmarshal(body, &obj)
	if err != nil {
		api.ReturnErr(api.GetBodyError, ctx, err, response)
		return
	}
	bodyString, err := json.Marshal(obj)
	if err != nil {
		api.ReturnErr(api.GetBodyError, ctx, err, response)
		return
	}

	job := Job{
		ID:         uuid.New().String(),
		CreateDate: time.Now(),
		JobInfo: JobInfo{
			TotalStage:  0,
			ActiveStage: 99,
			User:        api.GetUserName(ctx),
			State:       tasks.StateStarted,
			Option:      string(bodyString),
		},
	}

	var (
		bodySteps []StepInfo
		steps     []Step
	)
	if err = json.Unmarshal(body, &bodySteps); err != nil {
		api.ReturnErr(api.JsonUnmarshalError, ctx, err, response)
		return
	}

	for _, step := range bodySteps {
		if !MachineryInstance.IsTaskRegistered(step.Tag) {
			api.ReturnErr(api.CannotFindTaskError, ctx, fmt.Errorf("cannot find task function: %s", step.Name), response)
			return
		}
		if step.Stage > job.TotalStage {
			job.TotalStage = step.Stage
		}
		if step.Stage < job.ActiveStage {
			job.ActiveStage = step.Stage
		}
		steps = append(steps, Step{
			ID:         uuid.New().String(),
			JobId:      job.ID,
			CreateDate: time.Now(),
			State:      tasks.StatePending,
			StepInfo:   step,
		})
	}

	if job.ActiveStage == 99 || job.TotalStage < job.ActiveStage {
		api.ReturnErr(api.GetTaskBodyError, ctx,
			fmt.Errorf("get step body error, cannot get active(%d) or total(%d) stage, %s",
				job.ActiveStage, job.TotalStage, string(body)), response)
		return
	}

	_, err = instance.AddItem(&job, 1)
	if err != nil {
		api.ReturnErr(api.AddDbError, ctx, err, response)
		return
	}

	_, err = instance.AddItem(&steps, int64(len(steps)))
	if err != nil {
		api.ReturnErr(api.AddDbError, ctx, err, response)
		return
	}

	err = createTask(job.ID, job.ActiveStage)
	if err != nil {
		api.ReturnErr(api.CreateTaskError, ctx, err, response)
		return
	}
	api.ResponseBody(ctx, response, job)
}
