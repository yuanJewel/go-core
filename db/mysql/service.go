package mysql

import (
	"context"
	"encoding/json"
	"github.com/yuanJewel/go-core/db"
	"github.com/yuanJewel/go-core/db/object"
	gologger "github.com/yuanJewel/go-core/logger"
	"gorm.io/gorm"
	"math"
	"reflect"
	"runtime"
)

func (m *Mysql) HealthCheck() error {
	d := m.dbConn.Exec("select 1")
	return d.Error
}

func (m *Mysql) GetTables() ([]string, error) {
	return m.dbConn.Migrator().GetTables()
}

func (m *Mysql) HasTable(tableName string) bool {
	return m.dbConn.Migrator().HasTable(tableName)
}

func (m *Mysql) WithContext(ctx context.Context) db.Service {
	if m.mysqlConfig.redisInstance == nil {
		return &Mysql{dbConn: m.dbConn.WithContext(ctx), mysqlConfig: m.mysqlConfig}
	}
	return &Mysql{dbConn: m.dbConn.WithContext(ctx), mysqlConfig: &mysqlConfig{
		maxSearchLimit: m.mysqlConfig.maxSearchLimit,
		offsetPages:    m.mysqlConfig.maxSearchLimit,
		redisInstance:  m.mysqlConfig.redisInstance.WithContext(ctx),
	}}
}

func (m *Mysql) Preload(query string, args ...interface{}) db.Service {
	return &Mysql{dbConn: m.dbConn.Preload(query, args), mysqlConfig: m.mysqlConfig}
}

func (m *Mysql) Joins(query string, args ...interface{}) db.Service {
	return &Mysql{dbConn: m.dbConn.Joins(query, args), mysqlConfig: m.mysqlConfig}
}

func (m *Mysql) Where(query string, args ...interface{}) db.Service {
	return &Mysql{dbConn: m.dbConn.Where(query, args), mysqlConfig: m.mysqlConfig}
}

func (m *Mysql) Search(query string, keys string) db.Service {
	var (
		_idsInt []int
		_idsStr []string
	)
	if keys == "*" {
		return m
	} else {
		if err := json.Unmarshal([]byte(keys), &_idsInt); err != nil {
			if err := json.Unmarshal([]byte(keys), &_idsStr); err != nil {
				returnMysql := &Mysql{dbConn: m.dbConn, mysqlConfig: m.mysqlConfig}
				_ = returnMysql.dbConn.AddError(err)
				return returnMysql
			} else {
				return &Mysql{dbConn: m.dbConn.Where(query, _idsStr), mysqlConfig: m.mysqlConfig}
			}
		} else {
			return &Mysql{dbConn: m.dbConn.Where(query, _idsInt), mysqlConfig: m.mysqlConfig}
		}

	}
}

func (m *Mysql) Order(order string) db.Service {
	return &Mysql{dbConn: m.dbConn.Order(order), mysqlConfig: m.mysqlConfig}
}

func (m *Mysql) OffsetPages(pages int) db.Service {
	return &Mysql{dbConn: m.dbConn, mysqlConfig: &mysqlConfig{
		maxSearchLimit: m.mysqlConfig.maxSearchLimit,
		offsetPages:    pages,
		redisInstance:  m.mysqlConfig.redisInstance,
	}}
}

func (m *Mysql) Limit(limit int) db.Service {
	return &Mysql{dbConn: m.dbConn, mysqlConfig: &mysqlConfig{
		maxSearchLimit: limit,
		offsetPages:    m.mysqlConfig.offsetPages,
		redisInstance:  m.mysqlConfig.redisInstance,
	}}
}

// AddItem input must be an interface type
func (m *Mysql) AddItem(item interface{}, affectRows ...int64) (*gorm.DB, error) {
	if err := checkInput(item, reflect.Struct, reflect.Slice); err != nil {
		return nil, err
	}
	transaction := m.dbConn.Begin()
	result := transaction.Create(item)
	if result.Error != nil {
		return result, result.Error
	}
	if err := affectedRowsIsError(result, affectRows...); result.Error == nil && err != nil {
		transaction.Rollback()
		return result, err
	}
	if err := transaction.Commit().Error; err != nil {
		return result, err
	}
	return result, m.deleteCache(result)
}

// UpdateItem input new must be an interface type
func (m *Mysql) UpdateItem(old interface{}, new interface{}, affectRows ...int64) (*gorm.DB, error) {
	if err := checkInput(new, reflect.Struct); err != nil {
		return nil, err
	}
	transaction := m.dbConn.Begin()
	result := transaction.Where(old).Updates(new)
	if result.Error != nil {
		return result, result.Error
	}
	if err := affectedRowsIsError(result, affectRows...); result.Error == nil && err != nil {
		transaction.Rollback()
		return result, err
	}
	if err := transaction.Commit().Error; err != nil {
		return result, err
	}
	return result, m.deleteCache(result)
}

// DeleteItem input must be an interface type
func (m *Mysql) DeleteItem(item interface{}, affectRows ...int64) (*gorm.DB, error) {
	if err := checkInput(item, reflect.Struct, reflect.Slice); err != nil {
		return nil, err
	}
	transaction := m.dbConn.Begin()
	result := transaction.Delete(item)
	if result.Error != nil {
		return result, result.Error
	}
	if err := affectedRowsIsError(result, affectRows...); result.Error == nil && err != nil {
		transaction.Rollback()
		return result, err
	}
	if err := transaction.Commit().Error; err != nil {
		return result, err
	}
	return result, m.deleteCache(result)
}

// GetItems input get must be an interface type
func (m *Mysql) GetItems(find interface{}, get interface{}) (int64, error) {
	var total int64 = 0
	if err := checkInput(get, reflect.Slice); err != nil {
		return total, err
	}
	templateSql := m.dbConn.Where(find).Model(get)
	if err := m.queryByCache(templateSql, &total, func(db *gorm.DB) *gorm.DB {
		return db.Count(&total)
	}); err != nil {
		return total, err
	}

	return int64(math.Ceil(float64(total) / float64(m.mysqlConfig.maxSearchLimit))),
		m.queryByCache(templateSql.Offset(m.mysqlConfig.offsetPages*m.mysqlConfig.maxSearchLimit).
			Limit(m.mysqlConfig.maxSearchLimit), get, func(db *gorm.DB) *gorm.DB {
			return db.Find(get)
		})
}

// GetAllItems input get must be an interface type
func (m *Mysql) GetAllItems(find interface{}, get interface{}) (int64, error) {
	var total int64 = 0
	if err := checkInput(get, reflect.Slice); err != nil {
		return total, err
	}
	templateSql := m.dbConn.Where(find).Model(get)
	if err := m.queryByCache(templateSql, &total, func(db *gorm.DB) *gorm.DB {
		return db.Count(&total)
	}); err != nil {
		return total, err
	}

	return total, m.queryByCache(templateSql, get, func(db *gorm.DB) *gorm.DB {
		return db.Find(get)
	})
}

// GetItem input get must be an interface type
func (m *Mysql) GetItem(find interface{}, get interface{}) (bool, error) {
	if err := checkInput(get, reflect.Struct); err != nil {
		return false, err
	}

	var row int64 = 0
	// Check if the query rule will only get one
	templateSql := m.dbConn.Model(get).Where(find)
	if err := m.queryByCache(templateSql, &row, func(db *gorm.DB) *gorm.DB {
		return db.Count(&row)
	}); err != nil {
		return false, err
	}
	if row > 1 {
		var funcName = "unknown_function"
		pc, _, _, ok := runtime.Caller(1)
		if ok {
			funcName = runtime.FuncForPC(pc).Name()
		}
		gologger.Log.Errorf("Method %s executes query table data. The query logic can obtain %d records, but the program only needs one record.",
			funcName, row)
		return false, object.SelectOverOneError
	}
	if row == 0 {
		return false, nil
	}
	return true, m.queryByCache(templateSql, get, func(db *gorm.DB) *gorm.DB {
		return db.First(get)
	})
}

// checkInput Determine whether the input data is compliant
func checkInput(input interface{}, kinds ...reflect.Kind) error {
	if reflect.ValueOf(input).Kind() != reflect.Ptr || !isInSlice(reflect.ValueOf(input).Elem().Kind(), kinds) {
		var (
			funcName   = "unknown_function"
			funcFile   = ""
			funcLine   = 0
			optionName = "unknown_option"
		)
		pc, _, _, ok := runtime.Caller(1)
		if ok {
			optionName = runtime.FuncForPC(pc).Name()
		}
		pc2, pc2File, pc2Line, ok := runtime.Caller(2)
		if ok {
			funcName = runtime.FuncForPC(pc2).Name()
			funcFile = pc2File
			funcLine = pc2Line
		}
		gologger.Log.WithField("function", funcName).WithField("callerFile", funcFile).
			WithField("callerLine", funcLine).Errorf("方法 %s 执行 %s , 传入的对象(%s)不合法(*%s)",
			funcName, optionName, reflect.ValueOf(input).String(), printKinds(kinds))
		return object.InputError
	}
	return nil
}

func isInSlice(kind reflect.Kind, kindSlice []reflect.Kind) bool {
	for _, item := range kindSlice {
		if kind == item {
			return true
		}
	}
	return false
}

func printKinds(kindSlice []reflect.Kind) (kind string) {
	for _, item := range kindSlice {
		kind += item.String() + ", "
	}
	return
}

// affectedRowsIsError Determine whether the number of rows affected by database operations is as expected
func affectedRowsIsError(db *gorm.DB, right ...int64) error {
	if len(right) == 1 && right[0] == -1 {
		return nil
	}
	for _, r := range right {
		if db.RowsAffected == r {
			return nil
		}
	}

	gologger.Log.Errorf("%d records have been modified in database [%s], which is not the expected %d. The modification content is: %v",
		db.RowsAffected, db.Statement.Table, right, db.Statement.Dest)
	return object.AffectedRowsError
}
