package mysql

import (
	"fmt"
	"github.com/SmartLyu/go-core/config"
	"github.com/SmartLyu/go-core/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"time"
)

type Mysql struct {
	DbConn *gorm.DB
}

type mysqlLogger struct{}

func (mysqlLogger) Printf(format string, args ...interface{}) {
	logger.Log.Errorf(format, args...)
}

func GetMysqlInstance(cmdbCfgData *config.DataSourceDetail) (*Mysql, error) {
	newLogger := gormlogger.New(
		mysqlLogger{},
		gormlogger.Config{
			SlowThreshold:             time.Second,      // Slow SQL threshold
			LogLevel:                  gormlogger.Error, // Log level
			IgnoreRecordNotFoundError: true,             // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,            // Disable color
		},
	)

	dbConnString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local", cmdbCfgData.User, cmdbCfgData.Password, cmdbCfgData.Host, cmdbCfgData.Port, cmdbCfgData.Db, cmdbCfgData.Charset)
	dbConn, err := gorm.Open(mysql.Open(dbConnString), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return nil, err
	}
	mysqlInstance := &Mysql{DbConn: dbConn}

	sqlDB, err := dbConn.DB()
	if err != nil {
		return nil, err
	}
	// Set the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(cmdbCfgData.IdleConnections)
	// Set the maximum number of open database connections.
	sqlDB.SetMaxOpenConns(cmdbCfgData.MaxConnections)
	// Set the maximum time that a connection can be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)

	logger.Log.Infof("连接数据库 %s ，空闲连接数 %d ， 最大连接数 %d", cmdbCfgData.Host, cmdbCfgData.IdleConnections, cmdbCfgData.MaxConnections)
	return mysqlInstance, nil
}
