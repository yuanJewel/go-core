package db

import (
	"github.com/SmartLyu/go-core/cmdb"
	"github.com/SmartLyu/go-core/config"
)

func SetupCmdb(cfg *config.DataSourceDetail) error {
	if err := cmdb.InitCmdb(cfg); err != nil {
		return err
	}
	if err := cmdb.Instance.Setup([]interface{}{
		&Project{},
	}); err != nil {
		return err
	}
	return nil
}
