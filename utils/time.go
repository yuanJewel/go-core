package utils

import (
	"strings"
	"time"
)

const (
	Layout1                 = "2006-01-02T15:04:05.000Z"
	UTC8Layout              = "2006.01.02 15:04:05-0700"
	YYYY_MM_DD_Layout       = "2006-01-02"
	YYYYmmddHHMMSSLayout    = "20060102150405"
	YYYY_mm_dd_HHMMSSLayout = "2006-01-02 15:04:05"
)

// TimeParseYYYYMMDD 2006-01-02 类型的 字符串转换成 time.TIme
func TimeParseYYYYMMDD(value string) (time.Time, error) {
	return time.Parse(YYYY_MM_DD_Layout, value)
}

// TimeparseyyyyMmDdHhmmsslayout 2006-01-02 15:04:05 类型的 字符串转换成 time.TIme
// 时间需要是Utc 时间
func TimeparseyyyyMmDdHhmmsslayout(value string) (time.Time, error) {
	return time.Parse(YYYY_mm_dd_HHMMSSLayout, value)
}

// TimeToTimestampUnix 2006-01-02 15:04:05 类型的 字符串转换成 timestamp， 从 1790
// timestamp 是指格林威治时间1970年01月01日00时00分00秒
// (北京时间1970年01月01日08时00分00秒)起至现在的总秒数。
// 时间需要是 UTC 时间
func TimeToTimestampUnix(value string) (int64, error) {
	t, err := TimeparseyyyyMmDdHhmmsslayout(value)
	if err != nil {
		return 0, nil
	}
	return t.Unix(), nil
}

// TimeToTimestampUnixMilli 2006-01-02 15:04:05 类型的 字符串转换成 timestamp
// 毫秒数
// 时间需要是 UTC 时间
func TimeToTimestampUnixMilli(value string) (int64, error) {
	t, err := TimeToTimestampUnix(value)
	if err != nil {
		return 0, nil
	}
	return t * 1000, nil
}

// NowUtc 当前 UTC 时间的字符串
// 2006.01.02 15:04:05-0700
func NowUtc() string {
	return time.Now().Format(UTC8Layout)
}

// TimestampToTime 时间戳转化为时间 time.Time
func TimestampToTime(value int64) time.Time {
	t := time.Unix(value, 0)
	return t
}

// CstTimeToTimestampUnixMilli 时间戳转化为 时间字符串
// 转化后的格式为 2006-01-02 15:04:05
// 2006-01-02 15:04:05 类型的 字符串转换成 timestamp
// 毫秒数
// 时间需要是 CTS 时间 , UTC+8的时间
func CstTimeToTimestampUnixMilli(value string) (int64, error) {
	t, err := CstTimeToTimestampUnix(value)
	if err != nil {
		return 0, nil
	}
	return t * 1000, nil
}

// CstTimeToTimestampUnix 2006-01-02 15:04:05 类型的 字符串转换成 timestamp， 从 1790
// timestamp 是指格林威治时间1970年01月01日00时00分00秒
// (北京时间1970年01月01日08时00分00秒)起至现在的总秒数。
// 时间需要是 UTC 时间
func CstTimeToTimestampUnix(value string) (int64, error) {
	t, err := TimeparseyyyyMmDdHhmmsslayout(value)
	if err != nil {
		return 0, nil
	}
	t = t.Add(-8 * time.Hour)
	return t.Unix(), nil
}

// HourToFullFmt 输入: 15:04
// 要求输出 : 2006-01-02T15:04:05.000Z
func HourToFullFmt(timeStr string) (string, error) {
	timeStr = time.Now().Format("2006-01-02T") + strings.Trim(timeStr, " ")
	t, err := time.Parse("2006-01-02T15:04", timeStr)
	if err != nil {
		return "", err
	}
	s := t.Format(Layout1)
	return s, nil
}
