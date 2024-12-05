package task

import (
	"fmt"
	"github.com/yuanJewel/go-core/task"
	"time"
)

var RegisteredTask = map[string]interface{}{
	"test success": testSuccess,
	"test error":   testError,
}

func testError(id string, data ...interface{}) (string, error) {
	// Non-idempotent tasks require additional protection locks to use this module
	t := task.LockTaskState(id)
	if t != nil {
		return t.Error(), nil
	}

	err := task.SetVariable(id, "test-str", fmt.Sprintf("%v", data[0]))
	if err != nil {
		return "", err
	}
	s := fmt.Sprintf("task(%s) start in %s, input is %v, variable is %s", id, time.Now().String(), data,
		task.GetVariable(id, "test-str"))
	fmt.Println("yuanTag Error " + s)
	time.Sleep(1 * time.Second)
	return s, fmt.Errorf("%v", data)
}

func testSuccess(id string, data ...interface{}) (string, error) {
	s := fmt.Sprintf("task(%s) start in %s, input is %v, variable is %s", id, time.Now().String(),
		data, task.GetVariable(id, "test-list"))
	err := task.AppendVariable(id, "test-list", fmt.Sprintf("%v", data[0]))
	if err != nil {
		return "", err
	}
	fmt.Println("yuanTag Success " + s)
	time.Sleep(3 * time.Second)
	return s, nil
}
