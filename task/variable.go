package task

import (
	"fmt"
	"github.com/yuanJewel/go-core/db/service"
	"github.com/yuanJewel/go-core/logger"
	"time"
)

var stepToJob = make(map[string]string)

func SetVariable(i, key, value string) {
	id := getJobId(i)
	redisInstance.Set(-1, redisKey(id, key), value)
}

func GetVariable(i, key string) string {
	id := getJobId(i)
	value, err := redisInstance.Get(redisKey(id, key))
	if err != nil {
		return ""
	}
	return string(value)
}

// AppendVariable 分隔符为 ; 所以传递值中不能有v
func AppendVariable(i, key, value string) {
	id := getJobId(i)
	lock := true
	for lock {
		var err error
		lock, err = redisInstance.Exists(redisKey(id, key) + "_lock")
		if err != nil {
			logger.Log.Errorf("redis for task var lock is error, please check %v", err)
		}
		time.Sleep(time.Millisecond)
	}
	redisInstance.Set(-1, redisKey(id, key)+"_lock", time.Now().String())
	defer redisInstance.Del(redisKey(id, key) + "_lock")

	exists, err := redisInstance.Exists(redisKey(id, key))
	if err != nil {
		logger.Log.Errorf("redis for task var lock is error, please check %v", err)
	}
	if !exists {
		SetVariable(i, key, value)
		return
	}
	v, err := redisInstance.Get(redisKey(id, key))
	if err != nil {
		logger.Log.Errorf("redis for task var lock is error, please check %v", err)
		return
	}
	redisInstance.Set(-1, redisKey(id, key), fmt.Sprintf("%s;%s", v, value))
}

func redisKey(id, key string) string {
	return fmt.Sprintf("task_%s_var_%s", id, key)
}

func getJobId(stepId string) string {
	if jobId, ok := stepToJob[stepId]; ok {
		return jobId
	}
	var step Step
	exists, err := service.Instance.GetItem(Step{ID: stepId}, &step)
	if err != nil {
		logger.Log.Errorf("get job id from step %s is error, please check %v", stepId, err)
		return "Error_unknown_step_" + stepId
	}
	if !exists {
		logger.Log.Errorf("cannot find step: %s", stepId)
		return "Error_unknown_step_" + stepId
	}
	stepToJob[stepId] = step.JobId
	return step.JobId
}
