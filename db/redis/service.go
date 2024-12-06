package redis

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/sirupsen/logrus"
	"github.com/yuanJewel/go-core/api"
	"github.com/yuanJewel/go-core/logger"
)

// Ping checks Redis connection status
func (s *Store) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	if cancel != nil {
		defer cancel()
	}

	if err := s.redisInstance.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("redis ping failed: %w", err)
	}
	return nil
}

// WithContext creates a Store instance with context
func (s *Store) WithContext(ctx context.Context) *Store {
	return &Store{
		ctx:           ctx,
		expiration:    s.expiration,
		redisInstance: s.redisInstance,
		timeout:       s.timeout,
		retryDelay:    s.retryDelay,
	}
}

// Expire sets key expiration
func (s *Store) Expire(key string, expiration time.Duration) error {
	ctx, cancel := s.getContext()
	if cancel != nil {
		defer cancel()
	}
	if expiration == 0 {
		expiration = s.expiration
	}

	return s.redisInstance.Expire(ctx, key, expiration).Err()
}

func (s *Store) TTL(key string) (time.Duration, error) {
	ctx, cancel := s.getContext()
	if cancel != nil {
		defer cancel()
	}

	return s.redisInstance.TTL(ctx, key).Result()
}

func (s *Store) expire(key string, expiration time.Duration) error {
	now, err := s.TTL(key)
	if err != nil {
		return err
	}
	if expiration == 0 {
		expiration = s.expiration
	}
	if now < expiration {
		return s.Expire(key, expiration)
	}
	return nil
}

// Del deletes a key
func (s *Store) Del(key string) error {
	ctx, cancel := s.getContext()
	if cancel != nil {
		defer cancel()
	}

	if err := s.redisInstance.Del(ctx, key).Err(); err != nil {
		s.Log().Errorf("failed to delete key %s: %v", key, err)
		return err
	}
	return nil
}

// Exists checks if key exists
func (s *Store) Exists(key string) (bool, error) {
	ctx, cancel := s.getContext()
	if cancel != nil {
		defer cancel()
	}

	value, err := s.redisInstance.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("error checking key %s existence: %w", key, err)
	}
	return value == 1, nil
}

// Type returns key type
func (s *Store) Type(key string) (string, error) {
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

	keyType, err := s.redisInstance.Type(ctx, key).Result()
	if err != nil {
		return "", fmt.Errorf("failed to get type for key %s: %w", key, err)
	}
	return keyType, nil
}

// checkType checks if key has expected type
func (s *Store) checkType(key, expectedType string) error {
	exists, err := s.Exists(key)
	if err != nil {
		return fmt.Errorf("error checking key existence: %w", err)
	}
	if !exists {
		return nil
	}

	actualType, err := s.Type(key)
	if err != nil {
		return fmt.Errorf("error getting key type: %w", err)
	}

	if actualType != expectedType {
		return fmt.Errorf("wrong type for key %s, expected %s but got %s",
			key, expectedType, actualType)
	}

	return nil
}

// Lock acquires distributed lock
func (s *Store) Lock(key string, expiration time.Duration) bool {
	ctx, cancel := s.getContext()
	if cancel != nil {
		defer cancel()
	}

	for i := int(s.timeout / s.retryDelay); i > 0; i-- {
		if ctx.Err() != nil {
			s.Log().Errorf("lock acquisition timeout for key %s: %v", key, ctx.Err())
			return false
		}

		success, err := s.redisInstance.SetNX(ctx, key, "1", expiration).Result()
		if err != nil {
			s.Log().Errorf("failed to acquire lock for key %s: %v", key, err)
			return false
		}
		if success {
			return true
		}

		if i > 1 {
			select {
			case <-ctx.Done():
				return false
			case <-time.After(s.retryDelay):
			}
		}
	}

	s.Log().Warnf("failed to acquire lock for key %s: timeout after %v", key, s.timeout)
	return false
}

// Unlock releases distributed lock
func (s *Store) Unlock(key string) bool {
	ctx, cancel := s.getContext()
	if cancel != nil {
		defer cancel()
	}

	result, err := s.redisInstance.Del(ctx, key).Result()
	if err != nil {
		s.Log().Errorf("failed to release lock for key %s: %v", key, err)
		return false
	}
	return result > 0
}

// Log returns logger instance
func (s *Store) Log() *logrus.Entry {
	traceId := "-"
	if s.ctx != nil {
		if ctx, ok := s.ctx.(iris.Context); ok {
			traceId = api.GetTraceId(ctx)
		}
	}

	pc, file, line, ok := runtime.Caller(1)
	funcName := "unknown_function"
	optionName := "unknown_option"

	if ok {
		optionName = runtime.FuncForPC(pc).Name()
	}

	pc2, pc2File, pc2Line, ok := runtime.Caller(2)
	if ok {
		funcName = runtime.FuncForPC(pc2).Name()
		file = pc2File
		line = pc2Line
	}

	return logger.Log.WithFields(logrus.Fields{
		"traceId":    traceId,
		"function":   funcName,
		"callerFile": file,
		"callerLine": line,
		"option":     optionName,
	})
}

// getContext returns context and cancel function
func (s *Store) getContext() (context.Context, context.CancelFunc) {
	if s.ctx != nil {
		return s.ctx, nil
	}
	return context.WithTimeout(context.Background(), s.timeout)
}
