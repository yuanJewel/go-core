package redis

import (
	"fmt"
	"time"
)

// Set sets key-value pair
func (s *Store) Set(key string, value interface{}, expiration time.Duration) error {
	ctx, cancel := s.getContext()
	if cancel != nil {
		defer cancel()
	}
	if expiration == 0 {
		expiration = s.expiration
	}
	if err := s.checkType(key, "string"); err != nil {
		return err
	}

	if err := s.redisInstance.Set(ctx, key, value, expiration).Err(); err != nil {
		s.Log().Errorf("failed to set key-value pair - key: %s, error: %v", key, err)
		return err
	}
	return nil
}

// Get retrieves value by key
func (s *Store) Get(key string) (string, error) {
	ctx, cancel := s.getContext()
	if cancel != nil {
		defer cancel()
	}

	exists, err := s.Exists(key)
	if err != nil {
		return "", fmt.Errorf("error checking key %s existence: %w", key, err)
	}
	if !exists {
		return "", fmt.Errorf("key %s does not exist", key)
	}

	value, err := s.redisInstance.Get(ctx, key).Result()
	if err != nil {
		return "", fmt.Errorf("failed to get value for key %s: %w", key, err)
	}
	return value, nil
}
