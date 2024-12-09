package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

var Instance *Store

type Store struct {
	ctx           context.Context
	compress      bool
	timeout       time.Duration
	expiration    time.Duration
	retryDelay    time.Duration
	redisInstance *redis.Client
}
