package mysql

import (
	"fmt"
	"github.com/yuanJewel/go-core/config"
	"github.com/yuanJewel/go-core/db/redis"
	gologger "github.com/yuanJewel/go-core/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"time"
)

type Mysql struct {
	dbConn      *gorm.DB
	mysqlConfig *mysqlConfig
}

type mysqlConfig struct {
	redisInstance  *redis.Store
	maxSearchLimit int
	offsetPages    int
}

func GetMysqlInstance(cfgData *config.DataSourceDetail) (*Mysql, error) {
	logLevel := gormlogger.Warn
	if gologger.Log.Logger.GetLevel() < 3 {
		logLevel = gormlogger.LogLevel(gologger.Log.Logger.GetLevel())
	} else if gologger.Log.Logger.GetLevel() > 4 {
		logLevel = gormlogger.Info
	}
	newLogger := &logger{
		log: gologger.Log,
		Config: gormlogger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logLevel,    // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,       // Disable color
		},
	}

	dbConnString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local", cfgData.User, cfgData.Password, cfgData.Host, cfgData.Port, cfgData.Db, cfgData.Charset)
	dbConn, err := gorm.Open(mysql.Open(dbConnString), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return nil, err
	}
	mysqlInstance := &Mysql{dbConn: dbConn, mysqlConfig: &mysqlConfig{
		redisInstance:  redis.Instance,
		maxSearchLimit: cfgData.MaxSearchLimit,
		offsetPages:    0,
	}}

	sqlDB, err := dbConn.DB()
	if err != nil {
		return nil, err
	}
	// Set the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(cfgData.IdleConnections)
	// Set the maximum number of open database connections.
	sqlDB.SetMaxOpenConns(cfgData.MaxConnections)
	// Set the maximum time that a connection can be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)

	gologger.Log.Infof("连接数据库 %s, 空闲连接数 %d, 最大连接数 %d", cfgData.Host, cfgData.IdleConnections, cfgData.MaxConnections)
	return mysqlInstance, nil
}
