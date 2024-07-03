package mysql

import (
	"github.com/bpcoder16/Water/env"
	"github.com/bpcoder16/Water/logit"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"strconv"
	"time"
)

func init() {
	loadMySQLConfig()
}

func InitMySQL() {
	connectMaster()
	setMasterConnectionPool()
}

func connectMaster() {
	dsn := config.Master.Username + ":" + config.Master.Password +
		"@tcp(" + config.Master.Host + ":" + strconv.Itoa(config.Master.Port) + ")/" + config.Master.Database +
		"?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: dsn, // DSN data source name
		//DefaultStringSize:         256,   // string 类型字段的默认长度
		//DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		//DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		//DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		//SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
	}), &gorm.Config{
		Logger: NewLogger(logit.GetGlobalHelper(), logger.Config{
			SlowThreshold: 200 * time.Millisecond, // Slow SQL threshold
			LogLevel: func() logger.LogLevel {
				if env.RunMode() == env.RunModeRelease {
					return logger.Warn
				}
				return logger.Info
			}(), // Log level
			IgnoreRecordNotFoundError: true,  // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      false, // Don't include params in the SQL log
			Colorful:                  false,
		}),
	})
	if err != nil {
		panic(dsn + ", failed to connect database: " + err.Error())
	}
	defaultMySQLGormDBMaster = db
}

func setMasterConnectionPool() {
	sqlDB, _ := defaultMySQLGormDBMaster.DB()

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(config.Master.MaxIdleConns)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(config.Master.MaxOpenConns)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)
}
