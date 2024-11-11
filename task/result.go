package task

import (
	"encoding/json"
	"fmt"
	"github.com/RichardKnop/machinery/v1/tasks"
	"github.com/yuanJewel/go-core/db/service"
	"github.com/yuanJewel/go-core/logger"
	"regexp"
	"strings"
	"time"
)

var DesensitizationKeyList = []string{"password", "passwd", "token", "auth"}

func resultToDb(id string, _ ...interface{}) error {
	step := Step{}

	exists, err := service.Instance.GetItem(Step{ID: id}, &step)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("cannot find step: %s", id)
	}
	if step.State != tasks.StatePending && step.State != tasks.StateStarted && step.State != "" {
		return taskRunOverOnce(id)
	}

	task, err := MachineryInstance.GetBackend().GetState(id)
	if err != nil {
		return err
	}
	_, err = service.Instance.UpdateItem(Step{ID: id}, &Step{
		State:      task.State,
		Result:     HumanReadableResults(task.Results),
		FinishTime: time.Now(),
	}, 1)
	if err != nil {
		return err
	}
	return nil
}

func errorToDb(error, id string, _ ...interface{}) error {
	step := Step{}
	if strings.HasPrefix(id, "finish-") && strings.Contains(id, "-job-") {
		idList := strings.Split(id, "-job-")
		id = idList[0]
		_, err := service.Instance.AddItem(&Step{
			ID:         id,
			JobId:      idList[1],
			CreateDate: time.Now(),
			StartTime:  time.Now(),
			StepInfo: StepInfo{
				Name:   "任务完结模块抛出异常",
				Tag:    "finish",
				Stage:  0,
				Option: "",
			},
		}, 1)
		if err != nil {
			return err
		}
	}
	exists, err := service.Instance.GetItem(Step{ID: id}, &step)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("cannot find step: %s", id)
	}
	if step.State == tasks.StateFailure {
		return taskRunOverOnce(id)
	}

	task, err := MachineryInstance.GetBackend().GetState(id)
	if err != nil {
		return err
	}

	_, err = service.Instance.UpdateItem(Step{ID: id}, &Step{
		State:      task.State,
		Result:     HumanReadableResults(task.Results),
		Error:      error,
		FinishTime: time.Now(),
	}, 1)
	if err != nil {
		return err
	}

	_, err = service.Instance.UpdateItem(Step{JobId: step.JobId, State: tasks.StatePending},
		&Step{
			State: StateAborted,
			Error: fmt.Sprintf("task %s has failed, terminate this task", id),
		}, -1)
	if err != nil {
		return err
	}

	_, err = service.Instance.UpdateItem(Job{ID: step.JobId}, &Job{
		JobInfo: JobInfo{
			State:      tasks.StateFailure,
			FinishTime: time.Now(),
		},
	}, 1)
	if err != nil {
		return err
	}
	logger.Log.Infof("task %s (%s) finished failed", id, task.TaskName)
	return nil
}

func taskRunOverOnce(id string) error {
	logger.Log.Errorf("step %s has run over once!", id)
	return nil
}

func HumanReadableResults(taskResults []*tasks.TaskResult) string {
	resultValues := make([]string, 0, 0)
	flagError := false
	for _, taskResult := range taskResults {
		if taskResult.Type == "string" {
			v, ok := taskResult.Value.(string)
			if !ok {
				flagError = true
				break
			}
			resultValues = append(resultValues, desensitizationResults(v))
		} else if taskResult.Type == "[]string" {
			v, ok := taskResult.Value.([]interface{})
			if !ok {
				flagError = true
				break
			}
			for _, _v := range v {
				__v, ok := _v.(string)
				if !ok {
					flagError = true
					break
				}
				resultValues = append(resultValues, desensitizationResults(__v))
			}
		} else {
			flagError = true
			break
		}
	}
	if flagError {
		v, err := tasks.ReflectTaskResults(taskResults)
		if err != nil {
			return "reflect task results error"
		}
		return tasks.HumanReadableResults(v)
	}
	output, err := json.Marshal(resultValues)
	if err != nil {
		return strings.Join(resultValues, "\n")
	}
	return string(output)
}

func desensitizationResults(input string) string {
	var v interface{}
	err := json.Unmarshal([]byte(input), &v)
	if err != nil {
		if output, err := json.Marshal(map[string]string{"output": input}); err == nil {
			input = string(output)
		}
	}
	for _, key := range DesensitizationKeyList {
		reg := regexp.MustCompile(fmt.Sprintf("\"%s\":\"[^\"]*\"", key))
		input = reg.ReplaceAllString(input, fmt.Sprintf("\"%s\":\"***\"", key))
	}
	return input
}
