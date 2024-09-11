package cmdb

//
//import (
//	"encoding/json"
//	"fmt"
//	"github.com/fatih/structs"
//	"github.com/kataras/iris/v12"
//	"gorm.io/gorm"
//	"reflect"
//	"runtime"
//	"time"
//	"github.com/SmartLyu/go-core/api"
//	"github.com/SmartLyu/go-core/db"
//	"github.com/SmartLyu/go-core/logger"
//	"github.com/SmartLyu/go-core/utils"
//)
//
//func InsertAssetRecordItem(ctx iris.Context, body interface{}, context string, dbs ...gorm.DB) {
//	var (
//		function_name = "unknown_function"
//		function_file = ""
//		function_line = 0
//		user          db.User
//		username      = api.GetUserName(ctx)
//		chageTable    = []db.TableAffect{}
//		tranceid      = ctx.Request().Header.Get("traceid")
//		actionList    = []string{"INSERT", "SELECT", "UPDATE", "DELETE"}
//	)
//
//	pc, _pc_file, _pc_line, ok := runtime.Caller(1)
//	if ok {
//		function_name = runtime.FuncForPC(pc).Name()
//		function_file = _pc_file
//		function_line = _pc_line
//	}
//	_user_exist, err := CmdbInstance.GetItem(db.User{Name: username}, &user)
//	if err != nil {
//		logger.Log.Logger.WithField("function", function_name).WithField("callerfile", function_file).
//			WithField("callerline", function_line).Errorln(err)
//		return
//	}
//	if !_user_exist {
//		logger.Log.Logger.WithField("function", function_name).WithField("callerfile", function_file).
//			WithField("callerline", function_line).Errorf("cannot find user %s", username)
//		return
//	}
//	_body, err := json.Marshal(&body)
//	if err != nil {
//		logger.Log.Logger.WithField("function", function_name).WithField("callerfile", function_file).
//			WithField("callerline", function_line).Errorf("cannot get body: %v", body)
//	}
//
//	if _, err := CmdbInstance.AddItem(&db.AssetRecord{
//		TranceId:     tranceid,
//		UpdateTime:   time.Now(),
//		UpdateUserId: user.ID,
//		Url:          ctx.Path(),
//		Body:         string(_body),
//		Context:      context,
//	}, 1); err != nil {
//		logger.Log.Logger.WithField("function", function_name).WithField("callerfile", function_file).
//			WithField("callerline", function_line).Errorln(err)
//	}
//
//	for _, d := range dbs {
//		var (
//			table            = d.Statement.Table
//			action           = "unknow"
//			_unknown_clauses = []string{}
//			primaryIds       = readModel(d.Statement.Model)
//		)
//
//		if !CmdbInstance.HasTable(table) {
//			logger.Log.Logger.WithField("function", function_name).WithField("callerfile", function_file).
//				WithField("callerline", function_line).Errorf("cannot find table %s", table)
//		}
//		for _action := range d.Statement.Clauses {
//			if utils.InSlice(_action, actionList) {
//				action = _action
//				break
//			}
//			_unknown_clauses = append(_unknown_clauses, _action)
//		}
//		if action == "unknow" {
//			logger.Log.Logger.WithField("function", function_name).WithField("callerfile", function_file).
//				WithField("callerline", function_line).Errorf("cannot get action: %v", _unknown_clauses)
//		}
//		for _, primaryId := range primaryIds {
//			chageTable = append(chageTable, db.TableAffect{
//				TranceId:     tranceid,
//				UpdateTime:   time.Now(),
//				UpdateUserId: user.ID,
//				Table:        table,
//				PrimaryId:    primaryId,
//				Action:       action,
//			})
//		}
//	}
//	if len(chageTable) > 0 {
//		if _, err := CmdbInstance.AddItem(&chageTable, int64(len(chageTable))); err != nil {
//			logger.Log.Logger.WithField("function", function_name).WithField("callerfile", function_file).
//				WithField("callerline", function_line).Errorln(err)
//		}
//	}
//}
//
//func readModel(model interface{}) []string {
//	value := reflect.ValueOf(model)
//	if value.Kind() == reflect.Ptr {
//		value = value.Elem()
//	}
//	if value.Kind() == reflect.Slice {
//		var ids = []string{}
//		for i := 0; i < value.Len(); i++ {
//			_value := value.Index(i)
//			ids = append(ids, getId(_value.Interface()))
//		}
//		return ids
//	} else {
//		return []string{getId(value.Interface())}
//	}
//}
//
//func getId(value interface{}) string {
//	if reflect.ValueOf(value).Kind() != reflect.Struct {
//		logger.Log.Errorf("type is not struct: %v", value)
//		return "0"
//	}
//	_models := structs.New(value).Map()
//	for k, v := range _models {
//		if k == "ID" {
//			return fmt.Sprintf("%v", v)
//		}
//	}
//	return ""
//}
