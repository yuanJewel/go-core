package db

import (
	"fmt"
	"github.com/SmartLyu/go-core/cmdb"
	"github.com/SmartLyu/go-core/config"
	"github.com/SmartLyu/go-core/logger"
	"github.com/SmartLyu/go-core/utils"
)

func SetupCmdb(cfg *config.DataSourceDetail) error {
	if err := cmdb.InitCmdb(cfg); err != nil {
		return err
	}
	if err := cmdb.CmdbInstance.Setup([]interface{}{
		&Project{},
	}); err != nil {
		return err
	}
	_log, err := utils.ReadFromFile(logger.GetLogFilename())
	if err != nil {
		return err
	}
	fmt.Println(_log)
	return nil
}
