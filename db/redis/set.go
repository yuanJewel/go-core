package redis

import (
	"fmt"
	"time"
)

// SAdd 向集合添加元素
func (s *Store) SAdd(key string, expiration time.Duration, members ...interface{}) error {
	ctx, cancel := s.getContext()
	if cancel != nil {
		defer cancel()
	}
	if err := s.checkType(key, "set"); err != nil {
		return err
	}

	if err := s.redisInstance.SAdd(ctx, key, members...).Err(); err != nil {
		return fmt.Errorf("failed to add members to set - key: %s, error: %w", key, err)
	}
	return s.expire(key, expiration)
}

// SMembers 获取集合中的所有元素
func (s *Store) SMembers(key string) ([]string, error) {
	ctx, cancel := s.getContext()
	if cancel != nil {
		defer cancel()
	}
	if err := s.checkType(key, "set"); err != nil {
		return nil, err
	}

	members, err := s.redisInstance.SMembers(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get set members - key: %s, error: %w", key, err)
	}
	return members, nil
}
