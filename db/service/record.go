package service

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/structs"
	"github.com/kataras/iris/v12"
	"github.com/sirupsen/logrus"
	"github.com/yuanJewel/go-core/api"
	"github.com/yuanJewel/go-core/db/object"
	"github.com/yuanJewel/go-core/logger"
	"github.com/yuanJewel/go-core/utils"
	"gorm.io/gorm"
	"os"
	"reflect"
	"runtime"
	"time"
)

func InsertAssetRecordItem(ctx iris.Context, body interface{}, context string, dbs ...gorm.DB) {
	if !isRecordData() {
		return
	}

	var (
		functionName = "unknown_function"
		functionFile = ""
		functionLine = 0
		user         object.User
		username     = api.GetUserName(ctx)
		tranceId     = api.GetTraceId(ctx)
	)

	pc, pcFile, pcLine, ok := runtime.Caller(1)
	if ok {
		functionName = runtime.FuncForPC(pc).Name()
		functionFile = pcFile
		functionLine = pcLine
	}
	userExist, err := Instance.GetItem(object.User{Name: username}, &user)
	if err != nil {
		logger.Log.Logger.WithField("traceId", tranceId).WithField("function", functionName).
			WithField("callerFile", functionFile).
			WithField("callerLine", functionLine).Errorln(err)
		return
	}
	if !userExist {
		logger.Log.Logger.WithField("traceId", tranceId).WithField("function", functionName).
			WithField("callerFile", functionFile).
			WithField("callerLine", functionLine).Errorf("cannot find user %s", username)
		return
	}
	_body, err := json.Marshal(&body)
	if err != nil {
		logger.Log.Logger.WithField("traceId", tranceId).WithField("function", functionName).
			WithField("callerFile", functionFile).
			WithField("callerLine", functionLine).Errorf("cannot get body: %v", body)
	}

	if _, err := Instance.AddItem(&object.AssetRecord{
		TranceId:     tranceId,
		UpdateTime:   time.Now(),
		UpdateUserId: user.ID,
		Url:          ctx.Path(),
		Method:       ctx.Method(),
		Body:         string(_body),
		Context:      context,
	}, 1); err != nil {
		logger.Log.Logger.WithField("traceId", tranceId).WithField("function", functionName).
			WithField("callerFile", functionFile).WithField("callerLine", functionLine).Errorln(err)
	}

	entry := logger.Log.Logger.WithField("traceId", tranceId).WithField("function", functionName).
		WithField("callerFile", functionFile).WithField("callerLine", functionLine)
	AffectedTable(entry, ctx, user.ID, dbs...)
}

func AffectedTable(entry *logrus.Entry, ctx iris.Context, userId int, dbs ...gorm.DB) {
	if !isRecordData() {
		return
	}

	var (
		changeTable = make([]object.TableAffect, 0)
		actionList  = []string{"INSERT", "SELECT", "UPDATE", "DELETE"}
		tranceId    = api.GetTraceId(ctx)
	)

	for _, d := range dbs {
		var (
			table          = d.Statement.Table
			action         = "unknown"
			unknownClauses = make([]string, 0)
			primaryIds     = readModel(d.Statement.Model)
		)

		if isCheckTableExists() && !Instance.HasTable(table) {
			entry.Errorf("cannot find table %s", table)
		}
		for _action := range d.Statement.Clauses {
			if utils.InSlice(_action, actionList) {
				action = _action
				break
			}
			unknownClauses = append(unknownClauses, _action)
		}
		if action == "unknown" {
			entry.Errorf("cannot get action: %v", unknownClauses)
		}
		for _, primaryId := range primaryIds {
			changeTable = append(changeTable, object.TableAffect{
				TranceId:     tranceId,
				UpdateTime:   time.Now(),
				UpdateUserId: userId,
				Table:        table,
				PrimaryId:    primaryId,
				Action:       action,
			})
		}
	}
	if len(changeTable) > 0 {
		if _, err := Instance.AddItem(&changeTable, int64(len(changeTable))); err != nil {
			entry.Errorln(err)
		}
	}
}

func readModel(model interface{}) []string {
	value := reflect.ValueOf(model)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	if value.Kind() == reflect.Slice {
		var ids = []string{}
		for i := 0; i < value.Len(); i++ {
			_value := value.Index(i)
			ids = append(ids, getId(_value.Interface()))
		}
		return ids
	} else {
		return []string{getId(value.Interface())}
	}
}

func getId(value interface{}) string {
	if reflect.ValueOf(value).Kind() != reflect.Struct {
		logger.Log.Errorf("type is not struct: %v", value)
		return "0"
	}
	_models := structs.New(value).Map()
	for k, v := range _models {
		if k == "ID" {
			return fmt.Sprintf("%v", v)
		}
	}
	return ""
}

func isRecordData() bool {
	style := os.Getenv("RECORD_DATA")
	if style != "false" {
		return true
	}
	return false
}

func isCheckTableExists() bool {
	style := os.Getenv("CHECK_TABLE_EXISTS")
	if style != "true" {
		return false
	}
	return true
}
