package router

import (
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	router := gin.Default()
	r := router.Group("/pc_home")
	{

	}
	return router
}
