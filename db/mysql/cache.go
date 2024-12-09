package mysql

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	gologger "github.com/yuanJewel/go-core/logger"
	"gorm.io/gorm"
	"strings"
	"time"
)

const (
	cacheTableSetPrefix  = "cache:table:%s"
	lockTableCachePrefix = "cache:lock:table:%s"
	identityQueryIsEmpty = "nil"
)

func (m *Mysql) queryByCache(db *gorm.DB, dest interface{}, fc func(*gorm.DB) *gorm.DB) error {
	dryRun := fc(db.Session(&gorm.Session{DryRun: true, SkipDefaultTransaction: true, Logger: &nilLogger{}}))
	sqlStr := fmt.Sprintf("%v %s", db.Statement.Preloads,
		db.Dialector.Explain(dryRun.Statement.SQL.String(), dryRun.Statement.Vars...))

	hash := sha256.Sum256([]byte(sqlStr))
	key := fmt.Sprintf("cache:%s", hex.EncodeToString(hash[:]))
	if err := m.getCache(dryRun, key, sqlStr, dest); err == nil {
		return nil
	}

	result := fc(db)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		m.addCache(db, key, sqlStr, identityQueryIsEmpty)
		return nil
	}
	if result.Error == nil {
		m.addCache(db, key, sqlStr, dest)
	}
	return result.Error
}

func (m *Mysql) getCache(db *gorm.DB, key, sqlStr string, dest interface{}) error {
	if m.mysqlConfig.redisInstance == nil {
		return fmt.Errorf("redis instance is nil")
	}
	if m.checkIsLockTable(db) {
		return fmt.Errorf("cache is lock: %s", sqlStr)
	}

	startTime := time.Now().UnixMicro()
	cache, err := m.mysqlConfig.redisInstance.Get(key)
	if err == nil {
		if cache == identityQueryIsEmpty {
			gologger.Log.Debugf("[%.3fms] cache hit(%s) but is nil: %s", float64(time.Now().UnixMicro()-startTime)/1e3, key, sqlStr)
			return nil
		}
		err = m.mysqlConfig.redisInstance.Unmarshal([]byte(cache), dest)
		if err == nil {
			gologger.Log.Debugf("[%.3fms] cache hit(%s): %s", float64(time.Now().UnixMicro()-startTime)/1e3, key, sqlStr)
			return nil
		}
	}
	return fmt.Errorf("cache miss: %s", sqlStr)
}

func (m *Mysql) addCache(db *gorm.DB, key, sqlStr string, dest interface{}) {
	if m.mysqlConfig.redisInstance == nil {
		return
	}
	if m.checkIsLockTable(db) {
		return
	}

	var data string
	startTime := time.Now().UnixMicro()
	if dest != identityQueryIsEmpty {
		v, err := m.mysqlConfig.redisInstance.Marshal(dest)
		if err == nil {
			data = string(v)
		} else {
			return
		}
	} else {
		data = identityQueryIsEmpty
	}
	err := m.mysqlConfig.redisInstance.Set(key, data, 0)
	if err != nil {
		gologger.Log.Errorf("failed to add cache(%s): %s", key, sqlStr)
		return
	}
	for _, table := range getAffectTable(db) {
		err = m.mysqlConfig.redisInstance.SAdd(fmt.Sprintf(cacheTableSetPrefix, table), 0, key)
		if err != nil {
			_ = m.mysqlConfig.redisInstance.Del(key)
			gologger.Log.Errorf("failed to add cache(%s): %s", key, sqlStr)
			return
		}
	}

	gologger.Log.Debugf("[%.3fms] success to add cache(%s): %s", float64(time.Now().UnixMicro()-startTime)/1e3, key, sqlStr)
}

func (m *Mysql) deleteCache(db *gorm.DB) error {
	table := db.Statement.Table
	if m.mysqlConfig.redisInstance == nil {
		return nil
	}
	if exist, err := m.mysqlConfig.redisInstance.Exists(fmt.Sprintf(lockTableCachePrefix, table)); err == nil && exist {
		return nil
	}

	if err := m.mysqlConfig.redisInstance.Set(fmt.Sprintf(lockTableCachePrefix, table), time.Now().String(), 0); err != nil {
		return err
	}
	defer func() {
		_ = m.mysqlConfig.redisInstance.Del(fmt.Sprintf(lockTableCachePrefix, table))
	}()

	keys, err := m.mysqlConfig.redisInstance.SMembers(fmt.Sprintf(cacheTableSetPrefix, table))
	if err != nil {
		return fmt.Errorf("failed to get cache: %v", err)
	}
	for _, k := range keys {
		err = m.mysqlConfig.redisInstance.Del(k)
		if err != nil {
			return fmt.Errorf("failed to delete cache(%s): %v", k, err)
		}
	}

	if err = m.mysqlConfig.redisInstance.Del(fmt.Sprintf(cacheTableSetPrefix, table)); err != nil {
		return fmt.Errorf("failed to delete cache set(%s): %v", fmt.Sprintf(cacheTableSetPrefix, table), err)
	}
	gologger.Log.Debugf("success to delete cache(%d): %s", len(keys), table)
	return nil
}

func (m *Mysql) checkIsLockTable(db *gorm.DB) bool {
	for _, t := range getAffectTable(db) {
		if exist, err := m.mysqlConfig.redisInstance.Exists(fmt.Sprintf(lockTableCachePrefix, t)); err == nil && exist {
			return true
		}
	}
	return false
}

func getAffectTable(db *gorm.DB) []string {
	tables := make([]string, 0)
	tmp := make(map[string]bool)
	tmp[db.Statement.Table] = true
	for _, j := range db.Statement.Joins {
		tmp[extractTableName(j.Name)] = true
	}
	for preload, _ := range db.Statement.Preloads {
		if db.Statement.Schema != nil {
			if relation, ok := db.Statement.Schema.Relationships.Relations[preload]; ok {
				tmp[relation.FieldSchema.Table] = true
			}
		}
	}
	for t := range tmp {
		tables = append(tables, t)
	}
	return tables
}

func extractTableName(joinClause string) string {
	words := strings.Fields(joinClause)
	for i, word := range words {
		word = strings.ToUpper(word)
		if word == "JOIN" && i+1 < len(words) {
			return strings.TrimSpace(strings.Split(words[i+1], " ")[0])
		}
	}
	return strings.TrimSpace(strings.Split(joinClause, " ")[0])
}
