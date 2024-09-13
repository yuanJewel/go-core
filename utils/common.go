package utils

import (
	b64 "encoding/base64"
	"github.com/google/uuid"
)

// GenerateUuid  生成8位由数字字母组成的字符串
func GenerateUuid() string {
	projectid := uuid.New().String()[:8]
	return projectid
}

// RemoveRepByMap slice去重
func RemoveRepByMap(slc []string) []string {
	result := []string{}         // 存放返回的不重复切片
	tempMap := map[string]byte{} // 存放不重复主键
	for _, e := range slc {
		l := len(tempMap)
		tempMap[e] = 0 // 当e存在于tempMap中时，再次添加是添加不进去的，因为key不允许重复
		// 如果上一行添加成功，那么长度发生变化且此时元素一定不重复
		if len(tempMap) != l { // 加入map后，map长度变化，则元素不重复
			result = append(result, e) // 当元素不重复时，将元素添加到切片result中
		}
	}
	return result
}

func CreateUUID() string {
	return uuid.New().String()
}

// Base64Encode base64 加密
func Base64Encode(src []byte) string {
	dst := b64.StdEncoding.EncodeToString(src)
	return dst
}

// Base64Decode base64 解密
func Base64Decode(scr string) ([]byte, error) {
	return b64.StdEncoding.DecodeString(scr)
}
