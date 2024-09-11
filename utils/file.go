package utils

import (
	"fmt"
	"io/ioutil"
	"os"
)

// 判断文件是否存在
func FileExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// 遍历文件夹获取所有文件
func GetAllFile(pathname string, s []string) ([]string, error) {
	rd, err := ioutil.ReadDir(pathname)
	if err != nil {
		fmt.Println("read dir fail:", err)
		return s, err
	}
	for _, fi := range rd {
		if fi.IsDir() {
			fullDir := pathname + "/" + fi.Name()
			s, err = GetAllFile(fullDir, s)
			if err != nil {
				fmt.Println("read dir fail:", err)
				return s, err
			}
		} else {
			fullName := pathname + "/" + fi.Name()
			s = append(s, fullName)
		}
	}
	return s, nil
}

// 写入文件
func WriteToFile(filename, content string) error {
	err := ioutil.WriteFile(filename, []byte(content), 0644)
	return err
}

func BitToMb(b int64) int64 {
	return b / 1024 / 1024
}

func WritelinesToFile(filename, content string, _append bool) error {
	if !_append {
		return WriteToFile(filename, content)
	}
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	defer func() {
		_ = file.Close()
	}()
	if err != nil {
		return err
	}
	if _, err = file.WriteString(content); err != nil {
		return nil
	}
	return nil
}

func ReadFromFile(filename string) (string, error) {
	content, err := ioutil.ReadFile(filename)
	return string(content), err
}
