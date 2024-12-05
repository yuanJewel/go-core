package redis

import (
	"fmt"
	"time"
)

// LPush 将元素添加到列表头
func (s *Store) LPush(key string, value interface{}, expiration time.Duration) error {
	ctx, cancel := s.getContext()
	if cancel != nil {
		defer cancel()
	}
	if err := s.checkType(key, "list"); err != nil {
		return err
	}

	lockKey := fmt.Sprintf("lock:%s", key)
	if !s.Lock(lockKey, 5*time.Second) {
		return fmt.Errorf("unable to acquire redis list lock: %s", key)
	}
	defer s.Unlock(lockKey)

	if err := s.redisInstance.LPush(ctx, key, value).Err(); err != nil {
		return fmt.Errorf("failed to add redis list element - key: %s, error: %w", key, err)
	}
	return s.Expire(key, expiration)
}

// RPush 将元素添加到列表
func (s *Store) RPush(key string, value interface{}, expiration time.Duration) error {
	ctx, cancel := s.getContext()
	if cancel != nil {
		defer cancel()
	}
	if err := s.checkType(key, "list"); err != nil {
		return err
	}

	lockKey := fmt.Sprintf("lock:%s", key)
	if !s.Lock(lockKey, 5*time.Second) {
		return fmt.Errorf("unable to acquire redis list lock: %s", key)
	}
	defer s.Unlock(lockKey)

	if err := s.redisInstance.RPush(ctx, key, value).Err(); err != nil {
		return fmt.Errorf("failed to add redis list element - key: %s, error: %w", key, err)
	}
	return s.Expire(key, expiration)
}

// RPop 从列表中删除并返回最后一个元素
func (s *Store) RPop(key string) (string, error) {
	ctx, cancel := s.getContext()
	if cancel != nil {
		defer cancel()
	}
	if err := s.checkType(key, "list"); err != nil {
		return "", err
	}

	lockKey := fmt.Sprintf("lock:%s", key)
	if !s.Lock(lockKey, 5*time.Second) {
		return "", fmt.Errorf("unable to acquire redis list lock: %s", key)
	}
	defer s.Unlock(lockKey)

	value, err := s.redisInstance.RPop(ctx, key).Result()
	if err != nil {
		return "", fmt.Errorf("failed to delete redis list element - key: %s, error: %w", key, err)
	}
	return value, nil
}

// LPop 从列表中删除并返回第一个元素
func (s *Store) LPop(key string) (string, error) {
	ctx, cancel := s.getContext()
	if cancel != nil {
		defer cancel()
	}
	if err := s.checkType(key, "list"); err != nil {
		return "", err
	}

	lockKey := fmt.Sprintf("lock:%s", key)
	if !s.Lock(lockKey, 5*time.Second) {
		return "", fmt.Errorf("unable to acquire redis list lock: %s", key)
	}
	defer s.Unlock(lockKey)

	value, err := s.redisInstance.LPop(ctx, key).Result()
	if err != nil {
		return "", fmt.Errorf("failed to pop redis list element - key: %s, error: %w", key, err)
	}
	return value, nil
}

// LRange 获取列表指定范围的元素
func (s *Store) LRange(key string, start, stop int64) ([]string, error) {
	ctx, cancel := s.getContext()
	if cancel != nil {
		defer cancel()
	}
	if err := s.checkType(key, "list"); err != nil {
		return nil, err
	}

	lockKey := fmt.Sprintf("lock:%s", key)
	if !s.Lock(lockKey, 5*time.Second) {
		return nil, fmt.Errorf("unable to acquire redis list lock: %s", key)
	}
	defer s.Unlock(lockKey)

	values, err := s.redisInstance.LRange(ctx, key, start, stop).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get redis list element - key: %s, from %d to %d, error: %w",
			key, start, stop, err)
	}
	return values, nil
}

// LAll 获取列表所有元素
func (s *Store) LAll(key string) ([]string, error) {
	l, err := s.LLen(key)
	if err != nil {
		return nil, err
	}

	return s.LRange(key, 0, l-1)
}

// LLen 获取列表长度
func (s *Store) LLen(key string) (int64, error) {
	ctx, cancel := s.getContext()
	if cancel != nil {
		defer cancel()
	}
	if err := s.checkType(key, "list"); err != nil {
		return 0, err
	}

	length, err := s.redisInstance.LLen(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to get redis list length - key: %s, error: %w", key, err)
	}
	return length, nil
}
