package config

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"os"
)

var MysqlConfigMap = map[string]map[string]string{
	"debug": {
		"test":    "root:root@tcp(127.0.0.1:3306)/test?charset=utf8&parseTime=True&loc=Local&timeout=10ms",
		"company": "chelun:chelun@tcp(10.10.1.23:3306)/chelun_home?charset=utf8&parseTime=True&loc=Local&timeout=10ms",
	},
}
var redisConfigMap = map[string]string{
	"debug":   "127.0.0.1:6379",
	"test":    "192.168.255.38:6381",
	"release": "172.17.0.7:6380",
}
var redisConfigPwdMap = map[string]string{
	"debug":   "",
	"test":    "",
	"release": "A?t8hMT2z",
}

var RedisObj = redis.NewClient(&redis.Options{
	Addr:     "172.17.0.7:6380",
	Password: "A?t8hMT2z", // no password set
	DB:       0,           // use default DB
})

func SetRedisObj() {
	envname := gin.Mode()
	addr, ok := redisConfigMap[envname]
	pwd, _ := redisConfigPwdMap[envname]
	if len(envname) > 1 && ok {
		RedisObj = redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: pwd,
			DB:       0,
		})
	}
}

var currentPath, _ = os.Getwd()

var ModelPath = currentPath + "/models/"
var RouterPath = currentPath + "/router/"
var ApiPath = currentPath + "/handle/"
