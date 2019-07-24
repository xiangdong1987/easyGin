package main

import (
	"easyGin/config"
	orm "easyGin/database"
	"easyGin/router"
)

func main() {
	defer func() {
		for _, value := range orm.DBList {
			value.Close()
		}
	}()
	defer config.RedisObj.Close()
	config.SetRedisObj()
	r := router.InitRouter()
	r.Run(":8082")
}
