package config

type Server struct {
	Port              int    `default:"8080" yaml:"port" json:"port" env:"server.port"`
	Swagger           bool   `default:"true" yaml:"swagger" json:"swagger" env:"server.swagger"`
	HttpTimeout       int    `default:"10" yaml:"httpTimeout" json:"httpTimeout" env:"server.httpTimeout"`
	LogLevel          string `default:"info" yaml:"logLevel" json:"logLevel" env:"server.logLevel"`
	DisableStartupLog bool   `default:"true" yaml:"disableStartupLog" json:"disableStartupLog" env:"server.disableStartupLog"`
}

type Auth struct {
	Key          string `required:"true" yaml:"key" json:"key" env:"server.key"`
	Timeout      int    `default:"60" yaml:"timeout" json:"timeout" env:"auth.timeout"`
	Refresh      int    `default:"40" yaml:"refresh" json:"refresh" env:"auth.refresh"`
	CryptoKey    string `required:"true" yaml:"cryptoKey" json:"cryptoKey" env:"auth.cryptoKey"`
	CryptoPrefix string `required:"true" yaml:"cryptoPrefix" json:"cryptoPrefix" env:"auth.cryptoPrefix"`
}

type Ldap struct {
	Enable           bool     `yaml:"enable" json:"enable" env:"ldap.enable"`
	Host             string   `yaml:"host" json:"host" env:"ldap.host"`
	Port             string   `yaml:"port" json:"port" env:"ldap.port"`
	BindDn           string   `yaml:"bind_dn" json:"bind_dn" env:"ldap.bind_dn"`
	BindPassword     string   `yaml:"bind_password" json:"bind_password" env:"ldap.bind_password"`
	SearchBaseDn     string   `yaml:"search_base_dn" json:"search_base_dn" env:"ldap.search_base_dn"`
	SearchAttributes []string `yaml:"search_attributes" json:"search_attributes" env:"ldap.search_attributes"`
}

type Db struct {
	Driver           string `default:"mysql" yaml:"driver" json:"driver" env:"db.driver"`
	Host             string `required:"true" yaml:"host" json:"host" env:"db.host"`
	Db               string `required:"true" yaml:"db" json:"db" env:"db.db"`
	User             string `required:"true" yaml:"user" json:"user" env:"db.user"`
	Password         string `required:"true" yaml:"password" json:"password" env:"db.password"`
	Charset          string `default:"utf8" yaml:"charset" json:"charset" env:"db.charset"`
	Port             int    `default:"3306" yaml:"port" json:"port" env:"db.port"`
	CheckTableExists bool   `default:"true" yaml:"check_table_exists" json:"check_table_exists" env:"db.check_table_exists"`
	IdleConnections  int    `default:"1" yaml:"idle_connections" json:"idle_connections" env:"db.idle_connections"`
	MaxConnections   int    `default:"1" yaml:"max_connections" json:"max_connections" env:"db.max_connections"`
	MaxSearchLimit   int    `default:"100" yaml:"max_search_limit" json:"max_search_limit" env:"db.max_search_limit"`
}

type Redis struct {
	Host       string `required:"true" yaml:"host" json:"host" env:"redis.host"`
	Port       string `default:"6379" yaml:"port" json:"port" env:"redis.port"`
	Db         int    `required:"true" yaml:"db" json:"db" env:"redis.db"`
	Password   string `required:"true" yaml:"password" json:"password" env:"redis.password"`
	PoolSize   int    `default:"100" yaml:"pool_size" json:"pool_size" env:"redis.pool_size"`
	Timeout    int    `default:"5" yaml:"timeout" json:"timeout" env:"redis.timeout"`
	Expiration int    `default:"600" yaml:"expiration" json:"expiration" env:"redis.expiration"`
	RetryDelay int    `default:"100" yaml:"retry_delay_ms" json:"retry_delay_ms" env:"redis.retry_delay_ms"`
	IsZip      bool   `default:"false" yaml:"is_zip" json:"is_zip" env:"redis.is_zip"`
}
