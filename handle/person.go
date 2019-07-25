package handle

import (
	"easyGin/models"
	"easyGin/tools"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AddPerson(c *gin.Context) {
	var person models.Person
	if c.ShouldBind(&person) == nil {
		fmt.Println(person)
		err := person.Insert()
		if err != nil {
			c.JSON(http.StatusOK, err)
			return
		} else {
			c.JSON(http.StatusOK, tools.GetResult(0, "Success", ""))
			return
		}
	}
	c.JSON(http.StatusOK, tools.GetResult(0, "Params is wrong", ""))
	return
}
func GetPerson(c *gin.Context) {
	id := c.Param("id")
	var person models.Person
	if id != "" {
		err := person.GetById(tools.String2Int(id))
		if err != nil {
			c.JSON(http.StatusOK, err)
			return
		} else {
			c.JSON(http.StatusOK, tools.GetResult(0, "Success", gin.H{
				"person": person,
			}))
			return
		}
	} else {
		page := c.DefaultQuery("page", "1")
		pageSize := c.DefaultQuery("page_size", "10")
		list, pageInfo, err := person.GetList(tools.String2Int(page), tools.String2Int(pageSize))
		if err != nil {
			c.JSON(http.StatusOK, err)
			return
		} else {
			c.JSON(http.StatusOK, tools.GetResult(0, "Success", gin.H{
				"list":      list,
				"page_info": pageInfo,
			}))
			return
		}
	}
}
func DeletePerson(c *gin.Context) {
	id := c.Param("id")
	var person models.Person
	err := person.DeleteById(tools.String2Int(id))
	if err != nil {
		c.JSON(http.StatusOK, err)
		return
	} else {
		c.JSON(http.StatusOK, tools.GetResult(0, "Success", ""))
		return
	}
}
func ModifyPerson(c *gin.Context) {
	var person models.Person
	if c.ShouldBind(&person) == nil {
		fmt.Println(person)
		err := person.ModifyById()
		if err != nil {
			c.JSON(http.StatusOK, err)
			return
		} else {
			c.JSON(http.StatusOK, tools.GetResult(0, "Success", ""))
			return
		}
	}
	c.JSON(http.StatusOK, tools.GetResult(0, "Params is wrong", ""))
	return
}
