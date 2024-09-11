package cmdb

import (
	"errors"
	"github.com/SmartLyu/go-core/mysql"
	"github.com/SmartLyu/go-core/pkg/config"
	"strings"
)

var CmdbInstance CmdbService

// InitCmdb Get the cmdb instance, the default is mysql
func InitCmdb(cmdbCfgData *config.DataSourceDetail) error {
	cmdbDriver := cmdbCfgData.Driver
	var err error
	switch strings.ToLower(cmdbDriver) {
	case "oracle":
		return nil
	default:
		CmdbInstance, err = mysql.GetMysqlInstance(cmdbCfgData)
		if err != nil {
			return err
		}
	}
	if CmdbInstance == nil {
		return errors.New("config about config cannot find right instance")
	}
	return nil
}
