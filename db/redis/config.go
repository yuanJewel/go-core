package redis

type Redis struct {
	Host       string `required:"true" yaml:"host" env:"redis.host"`
	Port       string `default:"6379" yaml:"port" env:"redis.port"`
	Db         int    `required:"true" yaml:"db" env:"redis.db"`
	Password   string `required:"true" yaml:"password" env:"redis.password"`
	PoolSize   int    `default:"100" yaml:"pool_size" env:"redis.pool_size"`
	Timeout    int    `default:"5" yaml:"timeout" env:"redis.timeout"`
	Expiration int    `default:"10" yaml:"expiration" env:"redis.expiration"`
}
