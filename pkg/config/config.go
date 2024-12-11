package config

import (
	"github.com/yuanJewel/go-core/config"
	"github.com/yuanJewel/go-core/task"
)

var (
	GlobalConfig AppConfig
)

type AppConfig struct {
	config.BasicConfig `yaml:",inline"`
	Task               task.Task `yaml:"task" json:"task"`
}

func LoadConfig(filename string) error {
	return config.LoadConfig(filename, &GlobalConfig)
}
