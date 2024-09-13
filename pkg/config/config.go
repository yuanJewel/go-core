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
	Ldap                    `yaml:"ldap"`
}

type Ldap struct {
	Enable           bool     `yaml:"enable" env:"ldap.enable"`
	Host             string   `yaml:"host" env:"ldap.host"`
	Port             string   `yaml:"port" env:"ldap.port"`
	BindDn           string   `yaml:"bind_dn" env:"ldap.bind_dn"`
	BindPassword     string   `yaml:"bind_password" env:"ldap.bind_password"`
	SearchBaseDn     string   `yaml:"search_base_dn" env:"ldap.search_base_dn"`
	SearchAttributes []string `yaml:"search_attributes" env:"ldap.search_attributes"`
}

func LoadConfig(cfgfileName string) error {
	err := configor.New(&configor.Config{ErrorOnUnmatchedKeys: true}).Load(&GlobalConfig, cfgfileName)
	if err != nil {
		return err
	}
	return nil
}
