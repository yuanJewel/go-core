package cmdb

import (
	"encoding/json"
	"github.com/SmartLyu/go-core/api"
	"github.com/SmartLyu/go-core/utils"
	"github.com/kataras/iris/v12"
	"gorm.io/gorm"
	"net/http"
)

// GetDbByIdsInfo 根据ids获取数据
// search 搜索条件 必须是在db中定义的struct对象
// object 作为传入传出的结果集 必须是在db中定义的struct对象的slice指针
func GetDbByIdsInfo(ctx iris.Context, search, object interface{}) {
	response := api.ResponseInit(ctx)
	ids := ctx.GetHeader("ids")
	_, err := Instance.GetItemsByIds(search, object, ids)
	if err != nil {
		api.ReturnErr(api.SelectDbError, ctx, err, response)
		return
	}
	api.ResponseBody(ctx, response, object)
}

// PostDbInfo 批量添加数据
// 要求body传入的结构和db中定义的struct对象一致
// db中需要定义两层结构，detail json字段引用数据层，外加一个ID字段，作为uuid
// object 作为传入传出的结果集 必须是在db中定义的struct对象的slice指针
// special 函数，可以修改bodyObject，如添加自定义字段
func PostDbInfo(ctx iris.Context, object interface{}, special func(*map[string]interface{})) {
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
		_bodyObject := make(map[string]interface{})
		_exist, err := Instance.GetItems(_info, object)
		if err != nil {
			api.ReturnErr(api.SelectDbError, ctx, err, response)
			return
		}
		if !_exist {
			special(&_bodyObject)
			_bodyObject["detail"] = _info
			needAddSlice = append(needAddSlice, _bodyObject)
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
	api.ResponseBody(ctx, response, bodyInfo)
}

// PutDbByIdInfo 修改指定
// 要求body传入的结构和db中定义的struct对象一致
// db中需要定义两层结构，detail json字段引用数据层，外加一个ID字段，作为uuid
// special 函数，可以修改bodyObject，如添加自定义字段
func PutDbByIdInfo(ctx iris.Context, object interface{}, special func(*map[string]interface{})) {
	response := api.ResponseInit(ctx)
	id := ctx.GetHeader("id")
	body, err := ctx.GetBody()
	if err != nil {
		api.ReturnErr(api.GetBodyError, ctx, err, response)
		return
	}
	var (
		bodyInfo   = make(map[string]interface{})
		updateInfo = make(map[string]interface{})
	)
	err = json.Unmarshal(body, &bodyInfo)
	if err != nil {
		api.ReturnErr(api.JsonUnmarshalError, ctx, err, response)
		return
	}
	updateInfo["detail"] = bodyInfo
	updateInfo["id"] = id
	special(&updateInfo)

	if err = utils.MapToStruct(updateInfo, object); err != nil {
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

// DeleteDbByIdInfo 根据id删除数据
// path 为根据Id查看数据的接口名
// object 作为传入传出的结果集 必须是在db中定义的struct对象的slice指针
func DeleteDbByIdInfo(ctx iris.Context, path string, object interface{}) {
	response := api.ResponseInit(ctx)
	ids := ctx.GetHeader("ids")
	code, reverserBody := api.ReverserUtil(ctx, http.MethodGet, path)
	if code != 200 {
		errResponse, err := api.UnmarshalResponse(reverserBody)
		if err != nil {
			api.ReturnErr(api.UnmarshalReponseError, ctx, err, response)
			return
		}
		api.ResponseBody(ctx, response, errResponse)
		return
	}
	type responseStruct struct {
		api.Response
		Data []map[string]interface{} `json:"data,omitempty"`
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
	body := map[string]string{"header": ids}
	go InsertAssetRecordItem(ctx, body, "", results...)
	api.ResponseBody(ctx, response, body)
}
