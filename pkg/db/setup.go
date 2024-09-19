package db

import (
	"github.com/SmartLyu/go-core/cmdb"
)

func SetupCmdb() error {
	if err := cmdb.Instance.Setup([]interface{}{
		&Project{},
	}); err != nil {
		return err
	}
	return nil
}
