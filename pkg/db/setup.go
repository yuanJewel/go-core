package db

import (
	"github.com/yuanJewel/go-core/db/service"
	"github.com/yuanJewel/go-core/task"
)

func SetupCmdb() error {
	if err := service.Instance.Setup([]interface{}{
		&Project{}, &task.Job{}, &task.Step{},
	}); err != nil {
		return err
	}
	return nil
}
