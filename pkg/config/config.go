package config

import (
	"github.com/SmartLyu/go-core/config"
	"github.com/jinzhu/configor"
)

var (
	GlobalConfig AppConfig
)

type AppConfig struct {
	ApiVersion              string `required:"true" yaml:"apiVersion" env:"apiVersion"`
	config.Server           `yaml:"server"`
	config.Auth             `yaml:"auth"`
	config.DataSourceDetail `yaml:"db"`
}

func LoadConfig(cfgfileName string) error {
	err := configor.New(&configor.Config{ErrorOnUnmatchedKeys: true}).Load(&GlobalConfig, cfgfileName)
	if err != nil {
		return err
	}
	return nil
}
