package task

import (
	"fmt"
	"time"
)

var RegisteredTask = map[string]interface{}{
	"test success": testSuccess,
	"test error":   testError,
}

func testError(data ...interface{}) (string, error) {
	s := fmt.Sprintf("task start in %s, input is %v", time.Now().String(), data)
	time.Sleep(1 * time.Second)
	return s, fmt.Errorf("%v", data)
}

func testSuccess(data ...interface{}) (string, error) {
	s := fmt.Sprintf("task start in %s, input is %v", time.Now().String(), data)
	time.Sleep(5 * time.Second)
	return s, nil
}
