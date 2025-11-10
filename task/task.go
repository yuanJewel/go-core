package task

import (
	"encoding/json"
	"fmt"
	"github.com/RichardKnop/machinery/v1/tasks"
	"github.com/google/uuid"
	"github.com/yuanJewel/go-core/db"
	"time"
)

func CreateTask(job *Job, dbInstance db.Service, body []byte, user string) error {
	var obj interface{}
	err := json.Unmarshal(body, &obj)
	if err != nil {
		return err
	}
	bodyString, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	*job = Job{
		ID:         uuid.New().String(),
		CreateDate: time.Now(),
		JobInfo: JobInfo{
			TotalStage:  0,
			ActiveStage: 99,
			User:        user,
			State:       tasks.StateStarted,
			Option:      string(bodyString),
		},
	}

	var (
		bodySteps []StepInfo
		steps     []Step
	)
	if err = json.Unmarshal(body, &bodySteps); err != nil {
		return err
	}

	for _, step := range bodySteps {
		if !machineryInstance.IsTaskRegistered(step.Tag) {
			return fmt.Errorf("cannot find task function: %s", step.Name)
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
		return fmt.Errorf("get step body error, cannot get active(%d) or total(%d) stage, %s",
			job.ActiveStage, job.TotalStage, string(body))
	}

	_, err = dbInstance.AddItem(job, 1)
	if err != nil {
		return err
	}

	_, err = dbInstance.AddItem(&steps, int64(len(steps)))
	if err != nil {
		return err
	}

	err = createTask(job.ID, job.ActiveStage)
	return err
}
