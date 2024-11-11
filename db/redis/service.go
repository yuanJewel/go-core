package redis

import (
	"context"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/yuanJewel/go-core/logger"
	"runtime"
	"time"
)

func (s *Store) Set(expiration time.Duration, key string, value interface{}) {
	s.Lock()
	defer s.Unlock()
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	if expiration < 0 {
		expiration = s.expiration
	}
	data, err := json.Marshal(value)
	if err != nil {
		s.Log().Warningln(err)
		return
	}
	err = s.redisInstance.Set(ctx, key, data, expiration).Err()
	if err != nil {
		s.Log().Errorln(err)
	}
}

func (s *Store) Get(key string) ([]byte, error) {
	s.RLock()
	defer s.RUnlock()
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	value, err := s.redisInstance.Get(ctx, key).Bytes()
	if err != nil {
		s.Log().Warningln(err)
		return nil, err
	}
	return value, nil
}

func (s *Store) Keys(prefix string) ([]string, error) {
	s.RLock()
	defer s.RUnlock()
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	value, err := s.redisInstance.Keys(ctx, prefix).Result()
	if err != nil {
		s.Log().Warningln(err)
		return nil, err
	}
	return value, nil
}

func (s *Store) Log() *logrus.Entry {
	var (
		funcName   = "unknown_function"
		funcFile   = ""
		funcLine   = 0
		optionName = "unknown_option"
	)
	pc, _, _, ok := runtime.Caller(1)
	if ok {
		optionName = runtime.FuncForPC(pc).Name()
	}
	pc2, pc2File, pc2Line, ok := runtime.Caller(2)
	if ok {
		funcName = runtime.FuncForPC(pc2).Name()
		funcFile = pc2File
		funcLine = pc2Line
	}
	entry := logger.Log.WithField("function", funcName).WithField("callerFile", funcFile).
		WithField("callerLine", funcLine).WithField("option", optionName)
	return entry
}
