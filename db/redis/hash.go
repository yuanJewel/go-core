package redis

import "fmt"

// HSet 设置哈希表字段的值
func (s *Store) HSet(key, field string, value interface{}) error {
	ctx, cancel := s.getContext()
	if cancel != nil {
		defer cancel()
	}
	if err := s.checkType(key, "hash"); err != nil {
		return err
	}

	if err := s.redisInstance.HSet(ctx, key, field, value).Err(); err != nil {
		return fmt.Errorf("failed to set hash field - key: %s, field: %s, error: %w", key, field, err)
	}
	return nil
}

// HGet 获取哈希表中指定字段的值
func (s *Store) HGet(key, field string) (string, error) {
	ctx, cancel := s.getContext()
	if cancel != nil {
		defer cancel()
	}
	if err := s.checkType(key, "hash"); err != nil {
		return "", err
	}

	value, err := s.redisInstance.HGet(ctx, key, field).Result()
	if err != nil {
		return "", fmt.Errorf("failed to get hash field - key: %s, field: %s, error: %w", key, field, err)
	}
	return value, nil
}

// HGetAll 获取哈希表中所有的字段和值
func (s *Store) HGetAll(key string) (map[string]string, error) {
	ctx, cancel := s.getContext()
	if cancel != nil {
		defer cancel()
	}
	if err := s.checkType(key, "hash"); err != nil {
		return nil, err
	}

	values, err := s.redisInstance.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get all hash fields - key: %s, error: %w", key, err)
	}
	return values, nil
}

// HDel 删除哈希表中的一个或多个字段
func (s *Store) HDel(key string, fields ...string) error {
	ctx, cancel := s.getContext()
	if cancel != nil {
		defer cancel()
	}
	if err := s.checkType(key, "hash"); err != nil {
		return err
	}

	if err := s.redisInstance.HDel(ctx, key, fields...).Err(); err != nil {
		return fmt.Errorf("failed to delete hash fields - key: %s, error: %w", key, err)
	}
	return nil
}

// HExists 检查哈希表中是否存在指定的字段
func (s *Store) HExists(key, field string) (bool, error) {
	ctx, cancel := s.getContext()
	if cancel != nil {
		defer cancel()
	}
	if err := s.checkType(key, "hash"); err != nil {
		return false, err
	}

	exists, err := s.redisInstance.HExists(ctx, key, field).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check hash field existence - key: %s, field: %s, error: %w", key, field, err)
	}
	return exists, nil
}
