package redis

import (
	"github.com/go-redis/redis/v8"
	"sync"
	"time"
)

var Instance *Store

type Store struct {
	sync.RWMutex
	timeout       time.Duration
	expiration    time.Duration
	redisInstance *redis.Client
}
