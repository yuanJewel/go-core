package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/SmartLyu/go-core/api"
	"github.com/SmartLyu/go-core/utils"
	"github.com/kataras/iris/v12"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

type responseStruct struct {
	api.Response
	Data []map[string]interface{} `json:"data,omitempty"`
}

func initIdsString(ids string) string {
	returnIds := ids
	if !strings.HasPrefix(ids, "[") {
		if ids == "" || ids == "*" {
			returnIds = "*"
		} else {
			returnIds = fmt.Sprintf("[\"%s\"]", ids)
		}
	}
	return returnIds
}

// GetDbInfoByIds 根据ids获取数据
// search 搜索条件 必须是在db中定义的struct对象
// object 作为传入传出的结果集 必须是在db中定义的struct对象的slice指针
func GetDbInfoByIds(ctx iris.Context, search, object interface{}) {
	response := api.ResponseInit(ctx)
	ids := initIdsString(ctx.GetHeader("ids"))
	body, err := ctx.GetBody()
	if err != nil {
		api.ReturnErr(api.GetBodyError, ctx, err, response)
		return
	}
	err = json.Unmarshal(body, &search)
	if err != nil && len(body) != 0 {
		api.ReturnErr(api.JsonUnmarshalError, ctx, err, response)
		return
	}
	_, err = Instance.GetItemsByIds(search, object, ids)
	if err != nil {
		api.ReturnErr(api.SelectDbError, ctx, err, response)
		return
	}
	api.ResponseBody(ctx, response, object)
}

// GetDbInfoByIdsAndKey 根据ids获取数据
// search 搜索条件 必须是在db中定义的struct对象
// object 作为传入传出的结果集 必须是在db中定义的struct对象的slice指针
// key 作为需要新增检索的header字段
func GetDbInfoByIdsAndKey(ctx iris.Context, search, object interface{}, key string) {
	response := api.ResponseInit(ctx)
	ids := initIdsString(ctx.GetHeader("ids"))
	keys := initIdsString(ctx.GetHeader(key))
	body, err := ctx.GetBody()
	if err != nil {
		api.ReturnErr(api.GetBodyError, ctx, err, response)
		return
	}
	err = json.Unmarshal(body, &search)
	if err != nil && len(body) != 0 {
		api.ReturnErr(api.JsonUnmarshalError, ctx, err, response)
		return
	}
	_, err = Instance.GetItemsByIdsAndSlices(search, object, ids, key, keys)
	if err != nil {
		api.ReturnErr(api.SelectDbError, ctx, err, response)
		return
	}
	api.ResponseBody(ctx, response, object)
}

// PostDbInfo 批量添加数据
// 要求body传入的结构和db中定义的struct对象一致
// object 作为传入传出的结果集 必须是在db中定义的struct对象的slice指针
// special 函数，可以修改bodyObject，如添加自定义字段
func PostDbInfo(ctx iris.Context, object interface{}, special func(*map[string]interface{}) error) {
	response := api.ResponseInit(ctx)
	body, err := ctx.GetBody()
	if err != nil {
		api.ReturnErr(api.GetBodyError, ctx, err, response)
		return
	}
	var (
		bodyInfo     []map[string]interface{}
		needAddSlice = make([]map[string]interface{}, 0)
		results      = []gorm.DB{}
	)
	err = json.Unmarshal(body, &bodyInfo)
	if err != nil {
		api.ReturnErr(api.JsonUnmarshalError, ctx, err, response)
		return
	}

	for _, _info := range bodyInfo {
		if err = special(&_info); err != nil {
			api.ReturnErr(api.SpecialReturnError, ctx, err, response)
			return
		}
		_exist, err := Instance.GetItems(_info, object)
		if err != nil {
			api.ReturnErr(api.SelectDbError, ctx, err, response)
			return
		}
		if !_exist {
			needAddSlice = append(needAddSlice, _info)
		}
	}

	if err = utils.MapSliceToStructSlice(needAddSlice, object); err != nil {
		api.ReturnErr(api.ReflectError, ctx, err, response)
		return
	}

	if len(needAddSlice) > 0 {
		_result, err := Instance.AddItem(object, int64(len(needAddSlice)))
		if err != nil {
			api.ReturnErr(api.AddDbError, ctx, err, response)
			return
		}
		results = append(results, *_result)
	}
	go InsertAssetRecordItem(ctx, bodyInfo, "", results...)
	api.ResponseBody(ctx, response, needAddSlice)
}

// PutDbInfoById 修改指定
// 要求body传入的结构和db中定义的struct对象一致
// special 函数，可以修改bodyObject，如添加自定义字段
func PutDbInfoById(ctx iris.Context, path string, object interface{}, special func(*map[string]interface{}) error) {
	response := api.ResponseInit(ctx)
	id := ctx.GetHeader("id")
	header := http.Header{}
	header.Set("ids", id)
	code, reverserBody := api.ReverserInfoUtil(ctx, response, header, api.NilBody, http.MethodGet, path)
	if code != 200 {
		return
	}

	body, err := ctx.GetBody()
	if err != nil {
		api.ReturnErr(api.GetBodyError, ctx, err, response)
		return
	}
	var (
		returnObject responseStruct
		bodyInfo     = make(map[string]interface{})
	)
	if err = json.Unmarshal(reverserBody, &returnObject); err != nil {
		api.ReturnErr(api.GetBodyError, ctx, err, response)
		return
	}
	if len(returnObject.Data) == 0 {
		api.ReturnErr(api.GetBodyError, ctx,
			errors.New("the data that needs to be modified cannot be obtained according to the key in the header"),
			response)
		return
	}

	err = json.Unmarshal(body, &bodyInfo)
	if err != nil {
		api.ReturnErr(api.JsonUnmarshalError, ctx, err, response)
		return
	}

	bodyInfo["id"] = id
	if err = special(&bodyInfo); err != nil {
		api.ReturnErr(api.SpecialReturnError, ctx, err, response)
		return
	}

	if err = utils.MapToStruct(bodyInfo, object); err != nil {
		api.ReturnErr(api.ReflectError, ctx, err, response)
		return
	}
	result, err := Instance.UpdateItem(map[string]string{"id": id}, object, 1)
	if err != nil {
		api.ReturnErr(api.UpdateDbError, ctx, err, response)
		return
	}
	go InsertAssetRecordItem(ctx, bodyInfo, "", *result)
	api.ResponseBody(ctx, response, object)
}

// DeleteDb 根据id删除数据
// path 为根据Id查看数据的接口名
// object 作为传入传出的结果集 必须是在db中定义的struct对象的slice指针
func DeleteDb(ctx iris.Context, path string, object interface{}) {
	response := api.ResponseInit(ctx)
	code, reverserBody := api.ReverserUtil(ctx, response, http.MethodGet, path)
	if code != 200 {
		return
	}
	var (
		returnObject responseStruct
		results      = make([]gorm.DB, 0)
	)
	if err := json.Unmarshal(reverserBody, &returnObject); err != nil {
		api.ReturnErr(api.GetBodyError, ctx, err, response)
		return
	}
	for _, _d := range returnObject.Data {
		if err := utils.MapToStruct(_d, object); err != nil {
			api.ReturnErr(api.ReflectError, ctx, err, response)
			return
		}
		_result, err := Instance.DeleteItem(object, 1)
		if err != nil {
			api.ReturnErr(api.DeleteDbError, ctx, err, response)
			return
		}
		results = append(results, *_result)
	}

	go InsertAssetRecordItem(ctx, returnObject, "", results...)
	api.ResponseBody(ctx, response, returnObject)
}
