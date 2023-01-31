package dbconfig

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"ohmygin/nacosconfig"
)

var DBConnect *gorm.DB
var err error

func MysqlConnect() {
	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
	dsn := nacosconfig.Config.Mysql.DSN
	DBConnect, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if nil != err {
		log.Fatal(err)
	}
}
