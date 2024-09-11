package config

import (
	"github.com/jinzhu/configor"
)

var (
	GlobalConfig AppConfig
)

type AppConfig struct {
	ApiVersion       string `required:"true" yaml:"apiVersion" env:"apiVersion"`
	Server           `yaml:"server"`
	Auth             `yaml:"auth"`
	DataSourceDetail `yaml:"db"`
}

type Auth struct {
	Key string `required:"true" yaml:"key" env:"server.key"`
}

type Server struct {
	Port    int  `default:"8080" yaml:"port" env:"server.port"`
	Swagger bool `default:"true" yaml:"swagger" env:"server.swagger"`
}

type DataSourceDetail struct {
	Driver          string `default:"mysql" yaml:"driver" env:"db.driver"`
	Host            string `required:"true" yaml:"host" env:"db.host"`
	Db              string `required:"true" yaml:"db" env:"db.db"`
	User            string `required:"true" yaml:"user" env:"db.user"`
	Password        string `required:"true" yaml:"password" env:"db.password"`
	Charset         string `default:"utf8" yaml:"charset" env:"db.charset"`
	Port            int    `default:"3306" yaml:"port" env:"db.port"`
	IdleConnections int    `default:"1" yaml:"idleconnections" env:"db.idleconnections"`
	MaxConnections  int    `default:"1" yaml:"maxconnections" env:"db.maxconnections"`
}

func LoadConfig(cfgfileName string) error {
	err := configor.New(&configor.Config{ErrorOnUnmatchedKeys: true}).Load(&GlobalConfig, cfgfileName)
	if err != nil {
		return err
	}
	return nil
}
