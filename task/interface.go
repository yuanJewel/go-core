package task

import (
	"fmt"
	"github.com/RichardKnop/machinery/v1/tasks"
	"github.com/yuanJewel/go-core/db/service"
	"github.com/yuanJewel/go-core/logger"
	"time"
)

type FinishInterface interface {
	Success(string)
	Error(string)
	Abort(string)
}

var finishObject FinishInterface

type FinishStruct struct{}

func (f *FinishStruct) Success(id string) {
	logger.Log.Infof("Task %s finish success !", id)
}

func (f *FinishStruct) Error(id string) {
	logger.Log.Errorf("Task %s finish error !", id)
}

func (f *FinishStruct) Abort(id string) {
	logger.Log.Warnf("Task %s finish abort !", id)
}

func finishAbort(id string) error {
	_, err := service.Instance.UpdateItem(Step{JobId: id, State: tasks.StatePending},
		&Step{
			State: StateAborted,
			Error: fmt.Sprintf("job %s has failed or aborted, terminate this step", id),
		}, -1)
	if err != nil {
		return err
	}
	finishObject.Abort(id)
	return nil
}

func finishError(id string) error {
	if err := finishAbort(id); err != nil {
		return err
	}

	_, err := service.Instance.UpdateItem(Job{ID: id}, &Job{
		JobInfo: JobInfo{
			State:      tasks.StateFailure,
			FinishTime: time.Now(),
		},
	}, 0, 1)
	if err != nil {
		return err
	}

	finishObject.Error(id)
	return nil
}
