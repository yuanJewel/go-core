package task

import (
	"fmt"
	"github.com/RichardKnop/machinery/v1/tasks"
	"github.com/google/uuid"
	"github.com/yuanJewel/go-core/db/service"
	"github.com/yuanJewel/go-core/logger"
	"time"
)

func createTask(job string, stage int) error {
	var (
		steps      []Step
		signatures []*tasks.Signature
	)
	logger.Log.Debugf("start create task in job(%s[%d])", job, stage)
	number, err := service.Instance.GetAllItems(Step{JobId: job, StepInfo: StepInfo{Stage: stage}}, &steps)
	if err != nil {
		return err
	}
	if number == 0 {
		return finishTask(job, stage)
	}

	for _, step := range steps {
		signatures = append(signatures, &tasks.Signature{
			UUID: step.ID,
			Name: step.Tag,
			Args: []tasks.Arg{{
				Type:  "string",
				Value: step.ID,
			}, {
				Type:  "string",
				Value: step.Option,
			}},
			IgnoreWhenTaskNotRegistered: true,
			OnSuccess:                   newSignature(step.ID, "success"),
			OnError:                     newSignature(step.ID, "error"),
			RetryCount:                  0,
		})
	}

	_, err = service.Instance.UpdateItem(Step{JobId: job, StepInfo: StepInfo{Stage: stage}}, &Step{
		State:     tasks.StateStarted,
		StartTime: time.Now(),
	}, int64(len(steps)))
	if err != nil {
		return err
	}
	group, err := tasks.NewGroup(signatures...)
	if err != nil {
		return err
	}

	finishId := fmt.Sprintf("finish_%s", uuid.New().String())
	chord, err := tasks.NewChord(group, &tasks.Signature{
		UUID: finishId,
		Name: "finish",
		Args: []tasks.Arg{{
			Type:  "string",
			Value: job,
		}, {
			Type:  "int",
			Value: stage,
		}},
		IgnoreWhenTaskNotRegistered: true,
		RetryCount:                  0,
		OnSuccess:                   []*tasks.Signature{},
		OnError: []*tasks.Signature{{
			UUID: finishId + "_error",
			Name: "error",
			Args: []tasks.Arg{
				{
					Type:  "string",
					Value: finishId + "-job-" + job,
				},
			},
			RetryCount: 0,
		}},
	})
	if err != nil {
		return err
	}
	_, err = machineryInstance.SendChord(chord, maxConcurrency())
	if err == nil {
		logger.Log.Debugf("create task(%v) successfully", steps)
	}
	return err
}

func finishTask(jobId string, stage int, _ ...interface{}) error {
	var (
		job   Job
		steps []Step
	)

	logger.Log.Debugf("start finish task in job(%s[%d])", jobId, stage)
	exists, err := service.Instance.GetItem(Job{ID: jobId}, &job)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("job(%s) is not exists", jobId)
	}
	if job.State == StateAborted {
		return finishAbort(jobId)
	}

	if stage >= job.TotalStage {
		_, err = service.Instance.UpdateItem(Job{ID: jobId}, &Job{
			JobInfo: JobInfo{
				State:       tasks.StateSuccess,
				ActiveStage: job.TotalStage,
				FinishTime:  time.Now(),
			},
		}, 1)
		if err != nil {
			return err
		}
		finishObject.Success(jobId)
		return nil
	}

	_, err = service.Instance.GetAllItems(Step{JobId: jobId, StepInfo: StepInfo{Stage: stage}}, &steps)
	if err != nil {
		return err
	}

	for _, s := range steps {
		if s.State != tasks.StateSuccess {
			if s.State != tasks.StateStarted {
				logger.Log.Debugf("task(%s) is not success in job(%s[%d]), %v", s.ID, jobId, stage, s)
				return finishError(jobId)
			}
			task, err := machineryInstance.GetBackend().GetState(s.ID)
			if err != nil {
				logger.Log.Debugf("get task(%s) state error: %v", s.ID, err)
				return err
			}
			if task.State != tasks.StateSuccess {
				logger.Log.Debugf("task(%s) is not success in job(%s[%d]), %v", s.ID, jobId, stage, s)
				return finishError(jobId)
			}
		}
	}

	_, err = service.Instance.UpdateItem(Job{ID: jobId}, &Job{
		JobInfo: JobInfo{
			ActiveStage: job.ActiveStage + 1,
		},
	}, 1)
	if err != nil {
		return err
	}
	logger.Log.Debugf("finish task in job(%s[%d]) successfully", jobId, stage)
	return createTask(jobId, stage+1)
}

func newSignature(id string, status string) []*tasks.Signature {
	if status != "success" && status != "error" {
		return nil
	}
	return []*tasks.Signature{{
		UUID: id + "_" + status,
		Name: status,
		Args: []tasks.Arg{
			{
				Type:  "string",
				Value: id,
			},
		},
		RetryCount: 0,
	}}
}
