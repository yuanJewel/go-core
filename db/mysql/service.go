package mysql

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/SmartLyu/go-core/db"
	"github.com/SmartLyu/go-core/db/object"
	"github.com/SmartLyu/go-core/logger"
	"gorm.io/gorm"
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

func (m *Mysql) Preload(query string, args ...interface{}) db.Service {
	return &Mysql{dbConn: m.dbConn.Preload(query, args)}
}

func (m *Mysql) Joins(query string, args ...interface{}) db.Service {
	return &Mysql{dbConn: m.dbConn.Joins(query, args)}
}

func (m *Mysql) Where(query string, args ...interface{}) db.Service {
	return &Mysql{dbConn: m.dbConn.Where(query, args)}
}

// AddItem input must be an interface type
func (m *Mysql) AddItem(item interface{}, affectRows int64) (*gorm.DB, error) {
	if err := checkInput(item, reflect.Struct, reflect.Slice); err != nil {
		return nil, err
	}
	transaction := m.dbConn.Begin()
	result := transaction.Create(item)
	if result.Error != nil {
		return result, result.Error
	}
	if err := affectedRowsIsError(result.RowsAffected, affectRows); result.Error == nil && err != nil {
		transaction.Rollback()
		return result, err
	}
	return result, transaction.Commit().Error
}

// UpdateItem input new must be an interface type
func (m *Mysql) UpdateItem(old interface{}, new interface{}, affectRows int64) (*gorm.DB, error) {
	if err := checkInput(new, reflect.Struct); err != nil {
		return nil, err
	}
	transaction := m.dbConn.Begin()
	result := transaction.Where(old).Updates(new)
	if result.Error != nil {
		return result, result.Error
	}
	if err := affectedRowsIsError(result.RowsAffected, affectRows); result.Error == nil && err != nil {
		transaction.Rollback()
		return result, err
	}
	return result, transaction.Commit().Error
}

// DeleteItem input must be an interface type
func (m *Mysql) DeleteItem(item interface{}, affectRows int64) (*gorm.DB, error) {
	if err := checkInput(item, reflect.Struct, reflect.Slice); err != nil {
		return nil, err
	}
	transaction := m.dbConn.Begin()
	result := transaction.Delete(item)
	if result.Error != nil {
		return result, result.Error
	}
	if err := affectedRowsIsError(result.RowsAffected, affectRows); result.Error == nil && err != nil {
		transaction.Rollback()
		return result, err
	}
	return result, transaction.Commit().Error
}

// GetItems input get must be an interface type
func (m *Mysql) GetItems(find interface{}, get interface{}) (bool, error) {
	if err := checkInput(get, reflect.Slice); err != nil {
		return false, err
	}
	result := m.dbConn.Where(find).Find(get)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, nil
	}
	return result.RowsAffected != 0, result.Error
}

// GetItemsOrder input get must be an interface type
func (m *Mysql) GetItemsOrder(find interface{}, get interface{}, order string) (bool, error) {
	if err := checkInput(get, reflect.Slice); err != nil {
		return false, err
	}
	result := m.dbConn.Order(order).Where(find).Find(get)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, nil
	}
	return result.RowsAffected != 0, result.Error
}

// GetItemsFromDataAndSlice input get must be an interface type
func (m *Mysql) GetItemsFromDataAndSlice(find interface{}, query string, get interface{}, args ...interface{}) (bool, error) {
	if err := checkInput(get, reflect.Slice); err != nil {
		return false, err
	}
	result := m.dbConn.Where(find).Where(query, args...).Find(get)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, nil
	}
	return result.RowsAffected != 0, result.Error
}

// GetItemsFromDataAndSliceOrder input get must be an interface type
func (m *Mysql) GetItemsFromDataAndSliceOrder(find interface{}, query, order string, get interface{}, args ...interface{}) (bool, error) {
	if err := checkInput(get, reflect.Slice); err != nil {
		return false, err
	}
	result := m.dbConn.Order(order).Where(find).Where(query, args...).Find(get)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, nil
	}
	return result.RowsAffected != 0, result.Error
}

// GetItem input get must be an interface type
func (m *Mysql) GetItem(find interface{}, get interface{}) (bool, error) {
	if err := checkInput(get, reflect.Struct); err != nil {
		return false, err
	}

	// Check if the query rule will only get one
	row := m.dbConn.Model(get).Where(find).Find(&[]struct{}{}).RowsAffected
	if row > 1 {
		var _func_name = "unkown_function"
		pc, _, _, ok := runtime.Caller(1)
		if ok {
			_func_name = runtime.FuncForPC(pc).Name()
		}
		logger.Log.Errorf("方法 %s 执行插入表数据, 查询逻辑可以获取%d条记录，但是程序只需要一条记录",
			_func_name, row)
		return false, object.SelectOverOneError
	}

	result := m.dbConn.Where(find).First(get)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, nil
	}
	return true, result.Error
}

// GetItemsByIds input get must be an interface type
func (m *Mysql) GetItemsByIds(find interface{}, get interface{}, ids string) (bool, error) {
	var (
		_isString = false
		_idsInt   []int
		_idsStr   []string
		_exist    bool
		err       error
	)
	if ids == "*" {
		_exist, err = m.GetItems(find, get)
	} else {
		if err := json.Unmarshal([]byte(ids), &_idsInt); err != nil {
			if err := json.Unmarshal([]byte(ids), &_idsStr); err != nil {
				return false, errors.New("input ids is not a list")
			} else {
				_isString = true
			}
		}
		if _isString {
			_exist, err = m.GetItemsFromDataAndSlice(find, "id IN ?", get, _idsStr)
		} else {
			_exist, err = m.GetItemsFromDataAndSlice(find, "id IN ?", get, _idsInt)
		}
	}
	if err != nil {
		return false, err
	}
	return _exist, nil
}

// GetItemsByIdsOrder input get must be an interface type
func (m *Mysql) GetItemsByIdsOrder(find interface{}, get interface{}, ids, order string) (bool, error) {
	var (
		_isString = false
		_idsInt   []int
		_idsStr   []string
		_exist    bool
		err       error
	)
	if ids == "*" {
		_exist, err = m.GetItemsOrder(find, get, order)
	} else {
		if err := json.Unmarshal([]byte(ids), &_idsInt); err != nil {
			if err := json.Unmarshal([]byte(ids), &_idsStr); err != nil {
				return false, errors.New("input ids is not a list")
			} else {
				_isString = true
			}
		}
		if _isString {
			_exist, err = m.GetItemsFromDataAndSliceOrder(find, "id IN ?", order, get, _idsStr)
		} else {
			_exist, err = m.GetItemsFromDataAndSliceOrder(find, "id IN ?", order, get, _idsInt)
		}
	}
	if err != nil {
		return false, err
	}
	return _exist, nil
}

func (m *Mysql) GetItemsByIdsAndSlices(find interface{}, get interface{}, ids, name, slice string) (bool, error) {
	var (
		_slices []string
		_ids    []string
		_exist  bool
		err     error
	)
	if slice == "*" {
		return m.GetItemsByIds(find, get, ids)
	} else {
		if err := json.Unmarshal([]byte(slice), &_slices); err != nil {
			return false, err
		}
		if ids == "*" {
			_exist, err = m.GetItemsFromDataAndSlice(find, fmt.Sprintf("%s IN ?", name), get, _slices)
		} else {
			if err := json.Unmarshal([]byte(ids), &_ids); err != nil {
				return false, err
			}
			_exist, err = m.GetItemsFromDataAndSlice(find, fmt.Sprintf("id IN ? AND %s IN ?", name), get, _ids, _slices)
		}
	}

	if err != nil {
		return false, err
	}
	return _exist, nil
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
		logger.Log.WithField("function", funcName).WithField("callerFile", funcFile).
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
func affectedRowsIsError(num, right int64) error {
	if right != -1 && num != right {
		logger.Log.Errorf("添加用户存在异常，数据库发生%d条记录的修改，不是期望的%d", num, right)
		return object.AffectedRowsError
	}
	return nil
}
