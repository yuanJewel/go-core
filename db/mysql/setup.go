package mysql

import (
	"github.com/yuanJewel/go-core/db/object"
	gologger "github.com/yuanJewel/go-core/logger"
	"reflect"
)

func (m *Mysql) Setup(models []interface{}) error {
	var (
		notExistTables = []interface{}{}
	)
	models = append(models, []interface{}{
		&object.AssetRecord{}, &object.TableAffect{}, &object.User{},
	}...)
	defer gologger.Log.Infof("成功更新数据库，当前库中存在 %d 个数据表", len(models))

	err := m.dbConn.AutoMigrate(models...)
	if err != nil {
		for _, model := range models {
			if !m.dbConn.Migrator().HasTable(model) {
				notExistTables = append(notExistTables, model)
			} else {
				// https://gorm.io/docs/migration.html#Migrator-Interface
				// gorm officially does not have the ability to alert table, it needs to be implemented independently
				// Currently only new fields are supported
				_db := reflect.TypeOf(model)
				for i := 0; i < _db.NumField(); i++ {
					_columnName := _db.Field(i).Name
					if !m.dbConn.Migrator().HasColumn(model, _columnName) {
						gologger.Log.Infof("修改表 %s 新增字段 %s", _db.Name(), _columnName)
						if err := m.dbConn.Migrator().AddColumn(model, _columnName); err != nil {
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
			gologger.Log.Infof("新增表 %s", reflect.TypeOf(_notExistTable).Name())
		}
		err := m.dbConn.Migrator().CreateTable(notExistTables...)
		if err == nil {
			return nil
		}
	}
	return err
}
