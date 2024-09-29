package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/SmartLyu/go-core/api"
	"github.com/SmartLyu/go-core/logger"
	"github.com/SmartLyu/go-core/utils"
	"github.com/kataras/iris/v12"
	"gorm.io/gorm"
	"net/http"
	"strconv"
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

func initPageString(pages, traceId string) int {
	page := 1
	if pages != "" {
		parsedPage, err := strconv.Atoi(pages)
		if err != nil {
			logger.Log.WithField("traceId", traceId).WithField("function", "initPageString").
				Warnf("Failed to parse page number: %s, error: %v", pages, err)
		} else {
			page = parsedPage
		}
	}
	if page < 1 {
		page = 1
	}
	return page
}

func GetDbInfoByIds(ctx iris.Context, search, object interface{}, key ...string) {
	GetDbInfo(ctx, search, object, append([]string{"id"}, key...), []string{"id"})
}

func GetDbInfoByIdsAndOrder(ctx iris.Context, search, object interface{}, order ...string) {
	GetDbInfo(ctx, search, object, []string{"id"}, order)
}

// GetDbInfo 根据ids获取数据
// 会从header中获取两个字段，id和page
// id为检索条件，为空默认为*，检索所有
// page为分页好，为空默认为1
// search 搜索条件 必须是在db中定义的struct对象
// object 作为传入传出的结果集 必须是在db中定义的struct对象的slice指针
// key 作为需要新增检索的header字段
// order 为排序的字段，如果需要倒叙排列格式为"key desc"
func GetDbInfo(ctx iris.Context, search, object interface{}, keys, orders []string) {
	response := api.ResponseInit(ctx)
	page := initPageString(ctx.GetHeader("page"), response.TraceId)
	instance := Instance
	for _, key := range keys {
		if key != "" {
			keyString := initIdsString(ctx.GetHeader(key))
			instance = instance.Search(fmt.Sprintf("%s IN ?", key), keyString)
		}
	}

	for _, order := range orders {
		if order != "" {
			instance = instance.Order(order)
		}
	}

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

	nums, err := instance.OffsetPages(page-1).GetItems(search, object)
	if err != nil {
		api.ReturnErr(api.SelectDbError, ctx, err, response)
		return
	}
	ctx.Header("page", strconv.Itoa(page))
	ctx.Header("total-pages", strconv.FormatInt(nums, 16))
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
		nums, err := Instance.GetItems(_info, object)
		if err != nil {
			api.ReturnErr(api.SelectDbError, ctx, err, response)
			return
		}
		if nums == 0 {
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
	header.Set("id", id)
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
