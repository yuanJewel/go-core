package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/yuanJewel/go-core/logger"
	"time"
)

func InitRedis(cfg *Redis) (err error) {
	Instance, err = GetRedisInstance(cfg)
	return
}

func GetRedisInstance(cfg *Redis) (*Store, error) {
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
		expiration: time.Duration(cfg.Expiration) * time.Minute,
		timeout:    time.Duration(cfg.Timeout) * time.Second,
	}

	ctx, cancel := context.WithTimeout(context.Background(), instance.expiration)
	defer cancel()

	_, err := instance.redisInstance.Ping(ctx).Result()
	if err != nil {
		logger.Log.Errorf("cannot connect redis(%s:%s), error: %v", cfg.Host, cfg.Port, err)
		return nil, err
	}
	return instance, nil
}
