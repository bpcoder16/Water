package mysql

import "gorm.io/gorm"

var defaultMySQLGormDBMaster *gorm.DB

func GetMasterDB() *gorm.DB {
	return defaultMySQLGormDBMaster
}
