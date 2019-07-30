package router

import (
	. "easyGin/handle"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/index", IndexApi)
	//Add router

	return router
}
