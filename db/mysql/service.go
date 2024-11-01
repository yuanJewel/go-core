package mysql

import (
	"context"
	"encoding/json"
	"errors"
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
	return &Mysql{dbConn: m.dbConn.WithContext(ctx), mysqlConfig: m.mysqlConfig}
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
	return &Mysql{dbConn: m.dbConn, mysqlConfig: mysqlConfig{
		maxSearchLimit: m.mysqlConfig.maxSearchLimit,
		offsetPages:    pages,
	}}
}

func (m *Mysql) Limit(limit int) db.Service {
	return &Mysql{dbConn: m.dbConn, mysqlConfig: mysqlConfig{
		maxSearchLimit: limit,
		offsetPages:    m.mysqlConfig.offsetPages,
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
	if err := affectedRowsIsError(result.RowsAffected, affectRows...); result.Error == nil && err != nil {
		transaction.Rollback()
		return result, err
	}
	return result, transaction.Commit().Error
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
	if err := affectedRowsIsError(result.RowsAffected, affectRows...); result.Error == nil && err != nil {
		transaction.Rollback()
		return result, err
	}
	return result, transaction.Commit().Error
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
	if err := affectedRowsIsError(result.RowsAffected, affectRows...); result.Error == nil && err != nil {
		transaction.Rollback()
		return result, err
	}
	return result, transaction.Commit().Error
}

// GetItems input get must be an interface type
func (m *Mysql) GetItems(find interface{}, get interface{}) (int64, error) {
	var total int64 = 0
	if err := checkInput(get, reflect.Slice); err != nil {
		return total, err
	}
	result := m.dbConn.Where(find).Model(get).Count(&total).
		Offset(m.mysqlConfig.offsetPages * m.mysqlConfig.maxSearchLimit).
		Limit(m.mysqlConfig.maxSearchLimit).Find(get)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return total, nil
	}
	return int64(math.Ceil(float64(total) / float64(m.mysqlConfig.maxSearchLimit))), result.Error
}

// GetItem input get must be an interface type
func (m *Mysql) GetItem(find interface{}, get interface{}) (bool, error) {
	if err := checkInput(get, reflect.Struct); err != nil {
		return false, err
	}

	// Check if the query rule will only get one
	row := m.dbConn.Model(get).Where(find).Limit(m.mysqlConfig.maxSearchLimit).Find(&[]struct{}{}).RowsAffected
	if row > 1 {
		var _func_name = "unkown_function"
		pc, _, _, ok := runtime.Caller(1)
		if ok {
			_func_name = runtime.FuncForPC(pc).Name()
		}
		gologger.Log.Errorf("方法 %s 执行查询表数据, 查询逻辑可以获取%d条记录，但是程序只需要一条记录",
			_func_name, row)
		return false, object.SelectOverOneError
	}

	result := m.dbConn.Where(find).First(get)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, nil
	}
	return true, result.Error
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
		pc2, pc2_file, pc2_line, ok := runtime.Caller(2)
		if ok {
			funcName = runtime.FuncForPC(pc2).Name()
			funcFile = pc2_file
			funcLine = pc2_line
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
func affectedRowsIsError(num int64, right ...int64) error {
	if len(right) == 1 && right[0] == -1 {
		return nil
	}
	for _, r := range right {
		if num == r {
			return nil
		}
	}
	gologger.Log.Errorf("添加用户存在异常，数据库发生%d条记录的修改，不是期望的%d", num, right)
	return object.AffectedRowsError
}
