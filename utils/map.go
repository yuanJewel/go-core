package utils

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func InSlice(key interface{}, slice interface{}) bool {
	switch key.(type) {
	case int:
		data, ok := slice.([]int)
		if !ok {
			return false
		}
		for _, e := range data {
			if e == key {
				return true
			}
		}
	case string:
		data, ok := slice.([]string)
		if !ok {
			return false
		}
		for _, e := range data {
			if e == key {
				return true
			}
		}
	default:
		return false
	}
	return false
}

// MapToStruct 将单个 map 转换为 struct
func MapToStruct(m map[string]interface{}, obj interface{}) error {
	// 获取 obj 的类型
	objValue := reflect.ValueOf(obj).Elem()
	objType := objValue.Type()

	// 遍历结构体的字段
	for i := 0; i < objType.NumField(); i++ {
		field := objType.Field(i)
		fieldValue := objValue.Field(i)

		// 如果字段是嵌套结构体
		if field.Anonymous {
			nestedStruct := reflect.New(fieldValue.Type()).Interface()
			if err := MapToStruct(m, nestedStruct); err != nil {
				return err
			}
			fieldValue.Set(reflect.ValueOf(nestedStruct).Elem())
			continue
		}

		// 获取字段的 json 标签
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" || jsonTag == "-" {
			continue
		}
		if strings.Contains(jsonTag, ",") {
			jsonTag = strings.Split(jsonTag, ",")[0]
		}

		// 如果 json 标签匹配 map 中的键
		if value, ok := m[jsonTag]; ok {
			// 检查是否是子 map，需要递归
			if fieldValue.Kind() == reflect.Struct && fieldValue.Type() != reflect.TypeOf(time.Time{}) {
				nestedStruct := reflect.New(fieldValue.Type()).Interface()
				nestedMap, ok := value.(map[string]interface{})
				if ok {
					if err := MapToStruct(nestedMap, nestedStruct); err != nil {
						return err
					}
					fieldValue.Set(reflect.ValueOf(nestedStruct).Elem())
				}
			} else {
				if reflect.TypeOf(value).Kind() != fieldValue.Kind() {
					errReturn := fmt.Errorf("input map %s is unsupported kind: %v", jsonTag, fieldValue.Kind())
					switch v := value.(type) {
					case string:
						switch fieldValue.Kind() {
						case reflect.Int:
							intValue, err := strconv.Atoi(v)
							if err != nil {
								return err
							}
							fieldValue.SetInt(int64(intValue))
						case reflect.Bool:
							boolValue, err := strconv.ParseBool(v)
							if err != nil {
								return err
							}
							fieldValue.SetBool(boolValue)
						case reflect.Float64:
							floatValue, err := strconv.ParseFloat(v, 64)
							if err != nil {
								return err
							}
							fieldValue.SetFloat(floatValue)
						case reflect.Struct:
							if fieldValue.Type() == reflect.TypeOf(time.Time{}) {
								parsedTime, err := time.Parse(time.RFC3339, v)
								if err != nil {
									return err
								}
								fieldValue.Set(reflect.ValueOf(parsedTime).Convert(fieldValue.Type()))
							} else {
								return errReturn
							}
						default:
							return errReturn
						}
					case float64:
						switch fieldValue.Kind() {
						case reflect.Int:
							fieldValue.SetInt(int64(v))
						case reflect.Float64:
							fieldValue.SetFloat(v)
						case reflect.Bool:
							fieldValue.SetBool(v != 0)
						default:
							return errReturn
						}
					case int:
						switch fieldValue.Kind() {
						case reflect.Int:
							fieldValue.SetInt(int64(v))
						case reflect.Float64:
							fieldValue.SetFloat(float64(v))
						case reflect.Bool:
							fieldValue.SetBool(v != 0)
						default:
							return errReturn
						}
					default:
						return errReturn
					}
				} else {
					fieldValue.Set(reflect.ValueOf(value).Convert(fieldValue.Type()))
				}
			}
		}
	}
	return nil
}

// MapSliceToStructSlice 将 []map 转换为 []struct，并直接修改传入的切片
func MapSliceToStructSlice(maps []map[string]interface{}, outSlice interface{}) error {
	// 检查传入是否是指针类型
	v := reflect.ValueOf(outSlice)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Slice {
		return fmt.Errorf("outSlice 必须是指向 slice 的指针")
	}

	// 获取切片元素类型
	sliceValue := v.Elem()
	structType := sliceValue.Type().Elem()

	// 遍历 map 切片，逐个转换为 struct
	for _, m := range maps {
		// 创建一个新的 struct 实例
		structInstance := reflect.New(structType).Interface()

		// 将 map 转换为 struct
		err := MapToStruct(m, structInstance)
		if err != nil {
			return err
		}

		// 将转换后的 struct 添加到 slice 中
		structValue := reflect.ValueOf(structInstance).Elem()
		sliceValue.Set(reflect.Append(sliceValue, structValue))
	}

	return nil
}
