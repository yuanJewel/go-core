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

func GetRegisteredTaskNames() []string {
	return machineryInstance.GetRegisteredTaskNames()
}

func LockError(id string) error {
	return fmt.Errorf("task-lock-%s", id)
}

// LockTaskState Non-idempotent tasks require additional protection locks to use this module
func LockTaskState(id string) error {
	lockId := fmt.Sprintf("task_%s_lock", id)
	lock, err := redisInstance.Exists(lockId)
	if err != nil {
		logger.Log.Errorf("redis for task is error, please check %v", err)
		return nil
	}
	if lock {
		logger.Log.Debugf("task(%s) has been locked", id)
		return LockError(id)
	}
	redisInstance.Set(12*time.Hour, lockId, time.Now().String())
	return nil
}

func resultToDb(id string, _ ...interface{}) error {
	task, err := machineryInstance.GetBackend().GetState(id)
	if err != nil {
		return err
	}
	for _, r := range task.Results {
		if r.Type == "string" {
			v, ok := r.Value.(string)
			if !ok {
				break
			}
			if v == LockError(id).Error() {
				logger.Log.Warnf("step %s has been locked", id)
				return nil
			}
		}
	}

	t := LockTaskState(id + "_end")
	if t != nil {
		return nil
	}
	step := Step{}

	exists, err := service.Instance.GetItem(Step{ID: id}, &step)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("cannot find step: %s", id)
	}
	if step.State == tasks.StateSuccess {
		logger.Log.Errorf("step %s has run over once", id)
		return nil
	}

	_, err = service.Instance.UpdateItem(Step{ID: id}, &Step{
		State:      task.State,
		Result:     HumanReadableResults(task.Results),
		FinishTime: time.Now(),
	}, 1)
	if err != nil {
		return err
	}
	logger.Log.Debugf("result success step %s successfully", id)
	return nil
}

func errorToDb(error, id string, _ ...interface{}) error {
	t := LockTaskState(id + "_end")
	if t != nil {
		return nil
	}
	step := Step{}

	if strings.HasPrefix(id, "finish-") && strings.Contains(id, "-job-") {
		idList := strings.Split(id, "-job-")
		id = idList[0]
		finishObject.Error(idList[1])
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
		logger.Log.Errorf("step %s has run over once", id)
		return nil
	}

	task, err := machineryInstance.GetBackend().GetState(id)
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

	logger.Log.Debugf("result error step %s successfully", id)
	return finishError(step.JobId)
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
