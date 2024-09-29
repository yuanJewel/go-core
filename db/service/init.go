package service

import (
	"errors"
	"github.com/yuanJewel/go-core/config"
	"github.com/yuanJewel/go-core/db"
	"github.com/yuanJewel/go-core/db/mysql"
	"strings"
)

var Instance db.Service

// InitDb Get the cmdb instance, the default is mysql
func InitDb(cfgData *config.DataSourceDetail) error {
	cmdbDriver := cfgData.Driver
	var err error
	switch strings.ToLower(cmdbDriver) {
	case "oracle":
		return nil
	default:
		Instance, err = mysql.GetMysqlInstance(cfgData)
		if err != nil {
			return err
		}
	}
	if Instance == nil {
		return errors.New("config about config cannot find right instance")
	}
	return nil
}
