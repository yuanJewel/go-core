package config

type Server struct {
	Port              int    `default:"8080" yaml:"port" env:"server.port"`
	Swagger           bool   `default:"true" yaml:"swagger" env:"server.swagger"`
	HttpTimeout       int    `default:"10" yaml:"httpTimeout" env:"server.httpTimeout"`
	LogLevel          string `default:"info" yaml:"logLevel" env:"server.logLevel"`
	DisableStartupLog bool   `default:"true" yaml:"disableStartupLog" env:"server.disableStartupLog"`
}

type Auth struct {
	Key          string `required:"true" yaml:"key" env:"server.key"`
	Timeout      int    `default:"60" yaml:"timeout" env:"auth.timeout"`
	Refresh      int    `default:"40" yaml:"refresh" env:"auth.refresh"`
	CryptoKey    string `required:"true" yaml:"cryptoKey" env:"auth.cryptoKey"`
	CryptoPrefix string `required:"true" yaml:"cryptoPrefix" env:"auth.cryptoPrefix"`
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

type DataSourceDetail struct {
	Driver           string `default:"mysql" yaml:"driver" env:"db.driver"`
	Host             string `required:"true" yaml:"host" env:"db.host"`
	Db               string `required:"true" yaml:"db" env:"db.db"`
	User             string `required:"true" yaml:"user" env:"db.user"`
	Password         string `required:"true" yaml:"password" env:"db.password"`
	Charset          string `default:"utf8" yaml:"charset" env:"db.charset"`
	Port             int    `default:"3306" yaml:"port" env:"db.port"`
	CheckTableExists bool   `default:"true" yaml:"check_table_exists" env:"db.check_table_exists"`
	IdleConnections  int    `default:"1" yaml:"idle_connections" env:"db.idle_connections"`
	MaxConnections   int    `default:"1" yaml:"max_connections" env:"db.max_connections"`
	MaxSearchLimit   int    `default:"100" yaml:"max_search_limit" env:"db.max_search_limit"`
}
