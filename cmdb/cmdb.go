package cmdb

import (
	"errors"
	"github.com/SmartLyu/go-core/config"
	"github.com/SmartLyu/go-core/mysql"
	"strings"
)

var Instance Service

// InitCmdb Get the cmdb instance, the default is mysql
func InitCmdb(cmdbCfgData *config.DataSourceDetail) error {
	cmdbDriver := cmdbCfgData.Driver
	var err error
	switch strings.ToLower(cmdbDriver) {
	case "oracle":
		return nil
	default:
		Instance, err = mysql.GetMysqlInstance(cmdbCfgData)
		if err != nil {
			return err
		}
	}
	if Instance == nil {
		return errors.New("config about config cannot find right instance")
	}
	return nil
}
