package main

import (
	"go-pc_home/config"
	orm "go-pc_home/database"
	"go-pc_home/router"
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
	r.Run(":8081")
}
