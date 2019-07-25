package router

import (
	. "easyGin/handle"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/index", IndexApi)

	router.POST("/person/", AddPerson)
	router.GET("/person/:id", GetPerson)
	router.GET("/persons", GetPerson)
	router.DELETE("/person/:id", DeletePerson)
	router.PUT("/person/:id", ModifyPerson)

	//Add router
	return router
}
