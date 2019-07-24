package handle

import (
	"easyGin/tools"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func IndexApi(c *gin.Context) {
	fmt.Println("Welcome to EasyGin!")
	c.JSON(http.StatusOK, tools.GetResult(0, "Welcome to EasyGin!", gin.H{}))
}
