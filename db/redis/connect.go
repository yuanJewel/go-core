package redis

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/yuanJewel/go-core/config"
	"github.com/yuanJewel/go-core/logger"
	"time"
)

func InitRedis(cfg *config.Redis) (err error) {
	Instance, err = GetRedisInstance(cfg)
	return
}

func GetRedisInstance(cfg *config.Redis) (*Store, error) {
	instance := &Store{
		redisInstance: redis.NewClient(&redis.Options{
			Addr:         fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
			Password:     cfg.Password,
			DB:           cfg.Db,
			PoolSize:     cfg.PoolSize,
			DialTimeout:  time.Duration(cfg.Timeout+2) * time.Second,
			ReadTimeout:  time.Duration(cfg.Timeout) * time.Second,
			WriteTimeout: time.Duration(cfg.Timeout) * time.Second,
		}),
		expiration: time.Duration(cfg.Expiration) * time.Second,
		retryDelay: time.Duration(cfg.RetryDelay) * time.Millisecond,
		timeout:    time.Duration(cfg.Timeout) * time.Second,
		ctx:        nil,
	}

	err := instance.Ping()
	if err != nil {
		logger.Log.Errorf("cannot connect redis(%s:%s), error: %v", cfg.Host, cfg.Port, err)
		return nil, err
	}
	return instance, nil
}
