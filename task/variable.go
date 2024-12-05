package task

import (
	"fmt"
	"github.com/yuanJewel/go-core/db/service"
	"github.com/yuanJewel/go-core/logger"
)

// SetVariable 设置需要传递的变量，可以直接覆盖填写，该类型为字符串
func SetVariable(i, key, value string) error {
	id, err := getJobId(i)
	if err != nil {
		return err
	}
	varKey := redisKey(id, key)
	defer needRecycleKey(i, varKey)
	return redisInstance.Set(varKey, value, varExpiration)
}

func GetVariable(i, key string) string {
	jobId, err := getJobId(i)
	if err != nil {
		return ""
	}
	id := redisKey(jobId, key)
	t, err := redisInstance.Type(id)
	if err != nil {
		return ""
	}
	switch t {
	case "string":
		value, err := redisInstance.Get(id)
		if err != nil {
			return ""
		}
		return value
	case "list":
		value, err := redisInstance.LAll(id)
		if err != nil {
			return ""
		}
		return fmt.Sprintf("%v", value)
	}
	return ""
}

// AppendVariable 创建或增加需要传递的变量，用于追加参数，该类型为列表
func AppendVariable(i, key, value string) error {
	id, err := getJobId(i)
	if err != nil {
		return err
	}
	varKey := redisKey(id, key)
	defer needRecycleKey(i, varKey)
	return redisInstance.RPush(varKey, value, varExpiration)
}

func redisKey(id, key string) string {
	return fmt.Sprintf("job:%s:var:%s", id, key)
}

func getJobId(stepId string) (string, error) {
	if jobId, ok := stepToJob.Load(stepId); ok {
		return jobId.(string), nil
	}
	var step Step
	exists, err := service.Instance.GetItem(Step{ID: stepId}, &step)
	if err != nil {
		err = fmt.Errorf("get job id from step %s is error, please check %v", stepId, err)
		logger.Log.Errorln(err)
		return "", err
	}
	if !exists {
		err = fmt.Errorf("cannot find step: %s", stepId)
		logger.Log.Errorln(err)
		return "", err
	}
	stepToJob.Store(stepId, step.JobId)
	return step.JobId, nil
}
