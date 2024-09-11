package mysql

import (
	"encoding/json"
	"errors"
	object "github.com/SmartLyu/go-core/db"
	"github.com/SmartLyu/go-core/logger"
	"gorm.io/gorm"
	"reflect"
	"runtime"
)

func (this *Mysql) HealthCheck() error {
	db := this.DbConn.Exec("select 1")
	return db.Error
}

func (this *Mysql) GetTables() ([]string, error) {
	return this.DbConn.Migrator().GetTables()
}

func (this *Mysql) HasTable(tablename string) bool {
	return this.DbConn.Migrator().HasTable(tablename)
}

// AddItem input must be an interface type
func (this *Mysql) AddItem(item interface{}, affectRows int64) (*gorm.DB, error) {
	if err := checkInput(item, reflect.Struct, reflect.Slice); err != nil {
		return nil, err
	}
	transaction := this.DbConn.Begin()
	result := transaction.Create(item)
	if err := affectedRowsIsError(result.RowsAffected, affectRows); result.Error == nil && err != nil {
		transaction.Rollback()
		return result, err
	}
	return result, transaction.Commit().Error
}

// UpdateItem input new must be an interface type
func (this *Mysql) UpdateItem(old interface{}, new interface{}, affectRows int64) (*gorm.DB, error) {
	if err := checkInput(new, reflect.Struct); err != nil {
		return nil, err
	}
	transaction := this.DbConn.Begin()
	result := transaction.Where(old).Updates(new)
	if err := affectedRowsIsError(result.RowsAffected, affectRows); result.Error == nil && err != nil {
		transaction.Rollback()
		return result, err
	}
	return result, transaction.Commit().Error
}

// DeleteItem input must be an interface type
func (this *Mysql) DeleteItem(item interface{}, affectRows int64) (*gorm.DB, error) {
	if err := checkInput(item, reflect.Struct, reflect.Slice); err != nil {
		return nil, err
	}
	transaction := this.DbConn.Begin()
	result := transaction.Delete(item)
	if err := affectedRowsIsError(result.RowsAffected, affectRows); result.Error == nil && err != nil {
		transaction.Rollback()
		return result, err
	}
	return result, transaction.Commit().Error
}

// GetItems input get must be an interface type
func (this *Mysql) GetItems(find interface{}, get interface{}) (bool, error) {
	if err := checkInput(get, reflect.Slice); err != nil {
		return false, err
	}
	result := this.DbConn.Where(find).Find(get)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, nil
	}
	return result.RowsAffected != 0, result.Error
}

// GetItemsOrder input get must be an interface type
func (this *Mysql) GetItemsOrder(find interface{}, get interface{}, order string) (bool, error) {
	if err := checkInput(get, reflect.Slice); err != nil {
		return false, err
	}
	result := this.DbConn.Order(order).Where(find).Find(get)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, nil
	}
	return result.RowsAffected != 0, result.Error
}

// GetItemsFromSlice input get must be an interface type
func (this *Mysql) GetItemsFromSlice(find interface{}, get interface{}, args ...interface{}) (bool, error) {
	if err := checkInput(get, reflect.Slice); err != nil {
		return false, err
	}
	result := this.DbConn.Where(find, args...).Find(get)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, nil
	}
	return result.RowsAffected != 0, result.Error
}

// GetItemsFromSliceOrder input get must be an interface type
func (this *Mysql) GetItemsFromSliceOrder(find interface{}, order string, get interface{}, args ...interface{}) (bool, error) {
	if err := checkInput(get, reflect.Slice); err != nil {
		return false, err
	}
	result := this.DbConn.Order(order).Where(find, args...).Find(get)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, nil
	}
	return result.RowsAffected != 0, result.Error
}

// GetItemsFromDataAndSlice input get must be an interface type
func (this *Mysql) GetItemsFromDataAndSlice(find interface{}, query string, get interface{}, args ...interface{}) (bool, error) {
	if err := checkInput(get, reflect.Slice); err != nil {
		return false, err
	}
	result := this.DbConn.Where(find).Where(query, args...).Find(get)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, nil
	}
	return result.RowsAffected != 0, result.Error
}

// GetItemsFromDataAndSliceOrder input get must be an interface type
func (this *Mysql) GetItemsFromDataAndSliceOrder(find interface{}, query, order string, get interface{}, args ...interface{}) (bool, error) {
	if err := checkInput(get, reflect.Slice); err != nil {
		return false, err
	}
	result := this.DbConn.Order(order).Where(find).Where(query, args...).Find(get)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, nil
	}
	return result.RowsAffected != 0, result.Error
}

// GetItem input get must be an interface type
func (this *Mysql) GetItem(find interface{}, get interface{}) (bool, error) {
	if err := checkInput(get, reflect.Struct); err != nil {
		return false, err
	}

	// Check if the query rule will only get one
	row := this.DbConn.Model(get).Where(find).Find(&[]struct{}{}).RowsAffected
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

	result := this.DbConn.Where(find).First(get)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, nil
	}
	return true, result.Error
}

// GetItemsByIds input get must be an interface type
func (this *Mysql) GetItemsByIds(find interface{}, get interface{}, ids string) (bool, error) {
	var (
		_isString = false
		_idsInt   []int
		_idsStr   []string
		_exist    bool
		err       error
	)
	if ids == "*" {
		_exist, err = this.GetItems(find, get)
	} else {
		if err := json.Unmarshal([]byte(ids), &_idsInt); err != nil {
			if err := json.Unmarshal([]byte(ids), &_idsStr); err != nil {
				return false, err
			} else {
				_isString = true
			}
		}
		if _isString {
			_exist, err = this.GetItemsFromDataAndSlice(find, "id IN ?", get, _idsStr)
		} else {
			_exist, err = this.GetItemsFromDataAndSlice(find, "id IN ?", get, _idsInt)
		}
	}
	if err != nil {
		return false, err
	}
	return _exist, nil
}

// GetItemsByIdsOrder input get must be an interface type
func (this *Mysql) GetItemsByIdsOrder(find interface{}, get interface{}, ids, order string) (bool, error) {
	var (
		_isString = false
		_idsInt   []int
		_idsStr   []string
		_exist    bool
		err       error
	)
	if ids == "*" {
		_exist, err = this.GetItemsOrder(find, get, order)
	} else {
		if err := json.Unmarshal([]byte(ids), &_idsInt); err != nil {
			if err := json.Unmarshal([]byte(ids), &_idsStr); err != nil {
				return false, err
			} else {
				_isString = true
			}
		}
		if _isString {
			_exist, err = this.GetItemsFromDataAndSliceOrder(find, "id IN ?", order, get, _idsStr)
		} else {
			_exist, err = this.GetItemsFromDataAndSliceOrder(find, "id IN ?", order, get, _idsInt)
		}
	}
	if err != nil {
		return false, err
	}
	return _exist, nil
}

// GetItemsByStackIds input get must be an interface type
func (this *Mysql) GetItemsByStackIds(find interface{}, get interface{}, ids string) (bool, error) {
	var (
		_idsStr []string
		_exist  bool
		err     error
	)
	if ids == "*" {
		_exist, err = this.GetItems(find, get)
	} else {
		if err := json.Unmarshal([]byte(ids), &_idsStr); err != nil {
			return false, err
		} else {
			_exist, err = this.GetItemsFromDataAndSlice(find, "stack_id IN ?", get, _idsStr)
		}
	}
	if err != nil {
		return false, err
	}
	return _exist, nil
}

// GetItemsByIdsAndClouds input get must be an interface type
func (this *Mysql) GetItemsByIdsAndClouds(find interface{}, get interface{}, ids, clouds string) (bool, error) {
	var (
		_clouds []string
		_ids    []string
		_exist  bool
		err     error
	)
	if clouds == "*" {
		return this.GetItemsByIds(find, get, ids)
	} else {
		if err := json.Unmarshal([]byte(clouds), &_clouds); err != nil {
			return false, err
		}
		if ids == "*" {
			_exist, err = this.GetItemsFromDataAndSlice(find, "cloud_id IN ?", get, _clouds)
		} else {
			if err := json.Unmarshal([]byte(ids), &_ids); err != nil {
				return false, err
			}
			_exist, err = this.GetItemsFromDataAndSlice(find, "id IN ? AND cloud_id IN ?", get, _ids, _clouds)
		}
	}

	if err != nil {
		return false, err
	}
	return _exist, nil
}

// GetItemsByIdsAndProjects input get must be an interface type
func (this *Mysql) GetItemsByIdsAndProjects(find interface{}, get interface{}, ids, projects string) (bool, error) {
	var (
		_projects []string
		_ids      []string
		_exist    bool
		err       error
	)
	if projects == "*" {
		return this.GetItemsByIds(find, get, ids)
	} else {
		if err := json.Unmarshal([]byte(projects), &_projects); err != nil {
			return false, err
		}
		if ids == "*" {
			_exist, err = this.GetItemsFromDataAndSlice(find, "project_id IN ?", get, _projects)
		} else {
			if err := json.Unmarshal([]byte(ids), &_ids); err != nil {
				return false, err
			}
			_exist, err = this.GetItemsFromDataAndSlice(find, "id IN ? AND project_id IN ?", get, _ids, _projects)
		}
	}

	if err != nil {
		return false, err
	}
	return _exist, nil
}

func (this *Mysql) GetProjectSubsystemStack(projects, subsystems, stacks string, get interface{}) error {
	if err := checkInput(get, reflect.Slice); err != nil {
		return err
	}
	var (
		query = ""
		_ids  = []interface{}{}
	)
	db := this.DbConn.Table("projects").
		Select("*, projects.name as name, stacks.id as stack_id, stacks.name as stack_name, stacks.describe as stack_describe").
		Joins("inner join subsystems on projects.id = subsystems.project_id").
		Joins("inner join subsystem_stacks on subsystems.id = subsystem_stacks.subsystem_id").
		Joins("inner join stacks on stacks.id = subsystem_stacks.stack_id")
	if projects != "*" {
		var __ids []string
		if err := json.Unmarshal([]byte(projects), &__ids); err != nil {
			return err
		}
		if len(_ids) > 0 {
			query += " AND "
		}
		_ids = append(_ids, __ids)
		query += "projects.id IN ?"
	}
	if subsystems != "*" {
		var __ids []string
		if err := json.Unmarshal([]byte(subsystems), &__ids); err != nil {
			return err
		}
		if len(_ids) > 0 {
			query += " AND "
		}
		_ids = append(_ids, __ids)
		query += "subsystems.id IN ?"
	}
	if stacks != "*" {
		var __ids []string
		if err := json.Unmarshal([]byte(stacks), &__ids); err != nil {
			return err
		}
		if len(_ids) > 0 {
			query += " AND "
		}
		_ids = append(_ids, __ids)
		query += "stacks.id IN ?"
	}
	if len(_ids) > 0 {
		db.Where(query, _ids...)
	}
	result := db.Find(get)
	return result.Error
}

func (this *Mysql) GetProjectSubsystem(projects, subsystems string, get interface{}) error {
	if err := checkInput(get, reflect.Slice); err != nil {
		return err
	}
	var (
		query = ""
		_ids  = []interface{}{}
	)
	db := this.DbConn.Table("projects").
		Select("*, projects.name as name, subsystems.id as subsystem_id").
		Joins("inner join subsystems on projects.id = subsystems.project_id")
	if projects != "*" {
		var __ids []string
		if err := json.Unmarshal([]byte(projects), &__ids); err != nil {
			return err
		}
		if len(_ids) > 0 {
			query += " AND "
		}
		_ids = append(_ids, __ids)
		query += "projects.id IN ?"
	}
	if subsystems != "*" {
		var __ids []string
		if err := json.Unmarshal([]byte(subsystems), &__ids); err != nil {
			return err
		}
		if len(_ids) > 0 {
			query += " AND "
		}
		_ids = append(_ids, __ids)
		query += "subsystems.id IN ?"
	}
	if len(_ids) > 0 {
		db.Where(query, _ids...)
	}
	result := db.Find(get)
	return result.Error
}

// checkInput Determine whether the input data is compliant
func checkInput(input interface{}, kinds ...reflect.Kind) error {
	if reflect.ValueOf(input).Kind() != reflect.Ptr || !isInSlice(reflect.ValueOf(input).Elem().Kind(), kinds) {
		var (
			_func_name   = "unkown_function"
			_func_file   = ""
			_func_line   = 0
			_option_name = "unkown_option"
		)
		pc, _, _, ok := runtime.Caller(1)
		if ok {
			_option_name = runtime.FuncForPC(pc).Name()
		}
		pc2, pc2_file, pc2_line, ok := runtime.Caller(2)
		if ok {
			_func_name = runtime.FuncForPC(pc2).Name()
			_func_file = pc2_file
			_func_line = pc2_line
		}
		logger.Log.WithField("function", _func_name).WithField("callerfile", _func_file).
			WithField("callerline", _func_line).Errorf("方法 %s 执行 %s , 传入的对象(%s)不合法(*%s)",
			_func_name, _option_name, reflect.ValueOf(input).String(), printKinds(kinds))
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
