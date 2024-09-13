package mysql

import (
	"github.com/SmartLyu/go-core/db"
	"github.com/SmartLyu/go-core/logger"
	"reflect"
)

func (m *Mysql) Setup(models []interface{}) error {
	var (
		notExistTables = []interface{}{}
	)
	models = append(models, []interface{}{
		&db.AssetRecord{}, &db.TableAffect{}, &db.User{},
	}...)
	defer logger.Log.Infof("成功更新数据库，当前库中存在 %d 个数据表", len(models))

	err := m.DbConn.AutoMigrate(models...)
	if err != nil {
		for _, modle := range models {
			if !m.DbConn.Migrator().HasTable(modle) {
				notExistTables = append(notExistTables, modle)
			} else {
				// https://gorm.io/docs/migration.html#Migrator-Interface
				// gorm officially does not have the ability to alert table, it needs to be implemented independently
				// Currently only new fields are supported
				_db := reflect.TypeOf(modle)
				for i := 0; i < _db.NumField(); i++ {
					_columnname := _db.Field(i).Name
					if !m.DbConn.Migrator().HasColumn(modle, _columnname) {
						logger.Log.Infof("修改表 %s 新增字段 %s", _db.Name(), _columnname)
						if err := m.DbConn.Migrator().AddColumn(modle, _columnname); err != nil {
							return err
						}
					}
				}
			}
		}
		if len(notExistTables) == 0 {
			return nil
		}

		// gorm/schema/naming.go Line: 38, Func: TableName
		// The created table name may add an 's' ending
		for _, _notExistTable := range notExistTables {
			logger.Log.Infof("新增表 %s", reflect.TypeOf(_notExistTable).Name())
		}
		err := m.DbConn.Migrator().CreateTable(notExistTables...)
		if err == nil {
			return nil
		}
	}
	return err
}
