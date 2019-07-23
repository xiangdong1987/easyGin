package database

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"go-pc_home/config"
	"go-pc_home/tools"
)

var Eloquent *gorm.DB
var DBList map[string]*gorm.DB

func init() {
	var err error
	env := gin.Mode()
	dbConfig := config.MysqlConfigMap[env]
	if len(dbConfig) != 0 {
		DBList = make(map[string]*gorm.DB, len(dbConfig))
		for key, value := range dbConfig {
			DBList[key], err = gorm.Open("mysql", value)
			//日志
			if gin.Mode() == gin.DebugMode || gin.Mode() == gin.TestMode {
				DBList[key].LogMode(true)
			}
			if DBList[key].Error != nil {
				fmt.Printf("database error %v", DBList[key].Error)
			}
		}
	}
	if err != nil {
		fmt.Printf("mysql connect error %v", err)
	}
}

func Database(database string) (*gorm.DB, error) {
	var err error
	if database, ok := DBList[database]; ok {
		return database, err
	} else {
		return nil, tools.ReturnError{}.Instance(tools.ErrorDBConfig)
	}
}
