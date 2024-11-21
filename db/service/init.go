package service

import (
	"errors"
	"github.com/yuanJewel/go-core/config"
	"github.com/yuanJewel/go-core/db"
	"github.com/yuanJewel/go-core/db/mysql"
	"strings"
)

var Instance db.Service

// InitDb global variables are used to connect to the database by default.
func InitDb(cfgData *config.DataSourceDetail) (err error) {
	Instance, err = GetDb(cfgData)
	return
}

// GetDb Get the DataBase instance, the default is mysql
func GetDb(cfgData *config.DataSourceDetail) (instance db.Service, err error) {
	cmdbDriver := cfgData.Driver
	switch strings.ToLower(cmdbDriver) {
	case "oracle":
		return
	default:
		instance, err = mysql.GetMysqlInstance(cfgData)
		if err != nil {
			return
		}
	}
	if instance == nil {
		return nil, errors.New("config about config cannot find right instance")
	}
	return
}
