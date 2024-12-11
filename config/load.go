package config

import (
	"fmt"
	"github.com/jinzhu/configor"
	"reflect"
)

type BasicConfig struct {
	ApiVersion string `required:"true" yaml:"apiVersion" json:"apiVersion" env:"apiVersion"`
	Server     Server `yaml:"server" json:"server"`
	Auth       Auth   `yaml:"auth" json:"auth"`
	Redis      Redis  `yaml:"redis" json:"redis"`
	Db         Db     `yaml:"db" json:"db"`
}

func LoadConfig(filename string, cfg interface{}) error {
	v := reflect.ValueOf(cfg)
	if v.Kind() != reflect.Ptr {
		return fmt.Errorf("cfg must be a pointer")
	}

	if v.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("cfg must be a pointer to struct")
	}

	err := configor.New(&configor.Config{ErrorOnUnmatchedKeys: true}).Load(cfg, filename)
	if err != nil {
		return err
	}
	return nil
}
