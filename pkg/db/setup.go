package db

import (
	"github.com/yuanJewel/go-core/db/service"
	"github.com/yuanJewel/go-core/task"
)

func SetupCmdb() error {
	if err := service.Instance.Setup([]interface{}{
		// 用于测试的表
		&Project{},
		// task 需要用到任务系统的话，需要初始化这两个表
		&task.Job{}, &task.Step{},
	}); err != nil {
		return err
	}
	return nil
}
