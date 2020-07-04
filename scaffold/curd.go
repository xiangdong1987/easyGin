package scaffold

import (
	"bytes"
	"log"
	"strings"
	"text/template"
)

type CurdTemplate struct {
	LowerName   string
	StructName  string
	FirstChar   string
	DataIndex   string
	PrimaryKey  string
	PackageName string
}

const curdTemplate = `func ({{.FirstChar}} *{{.StructName}}) Insert() (err error) {
	dbe, err := database.Database("{{.DataIndex}}") 
	if err != nil { 
		return 
	} 
	err = dbe.Create({{.FirstChar}}).Error 
	return 
}
func ({{.FirstChar}} *{{.StructName}}) GetById(id int) (err error) {
	dbe, err := database.Database("{{.DataIndex}}")
	if err != nil {
		return
	}
	err = dbe.Where("{{.PrimaryKey}}=?", id).First({{.FirstChar}}).Error
	return
}
func ({{.FirstChar}} *{{.StructName}}) GetList(currentPage int, pageSize int) (list []{{.StructName}}, pageInfo database.PageInfo, err error) { 
	dbe, err := database.Database("{{.DataIndex}}")
	if err != nil {
		return
	}
	start := (currentPage - 1) * pageSize
	dbList := dbe.Limit(pageSize).Offset(start).Order("id DESC")
	err = dbList.Find(&list).Error
	pageInfo = pageInfo.GetPageInfo(dbe, &{{.StructName}}{}, currentPage, pageSize)
	return
}
func ({{.FirstChar}} *{{.StructName}}) DeleteById(id int) (err error) {
	dbe, err := database.Database("{{.DataIndex}}")
	if err != nil {
		return
	}
	err = dbe.Where("{{.PrimaryKey}}=?", id).Delete({{.FirstChar}}).Error
	return
}
func ({{.FirstChar}} *{{.StructName}}) ModifyById() (err error) {
	dbe, err := database.Database("{{.DataIndex}}")
	if err != nil {
		return
	}
	err = dbe.Save({{.FirstChar}}).Error
	return
}`

const apiTemplate = `package handle

import (
	"{{.PackageName}}/models"
	"{{.PackageName}}/tools"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Add{{.StructName}}(c *gin.Context) {
	var {{.LowerName}} models.{{.StructName}}
	if c.ShouldBind(&{{.LowerName}}) == nil {
		fmt.Println({{.LowerName}})
		err := {{.LowerName}}.Insert()
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
func Get{{.StructName}}(c *gin.Context) {
	id := c.Param("id")
	var {{.LowerName}} models.{{.StructName}}
	if id != "" {
		err := {{.LowerName}}.GetById(tools.String2Int(id))
		if err != nil {
			c.JSON(http.StatusOK, err)
			return
		} else {
			c.JSON(http.StatusOK, tools.GetResult(0, "Success", gin.H{
				"{{.LowerName}}": {{.LowerName}},
			}))
			return
		}
	} else {
		page := c.DefaultQuery("page", "1")
		pageSize := c.DefaultQuery("page_size", "10")
		list, pageInfo, err := {{.LowerName}}.GetList(tools.String2Int(page), tools.String2Int(pageSize))
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
func Delete{{.StructName}}(c *gin.Context) {
	id := c.Param("id")
	var {{.LowerName}} models.{{.StructName}}
	err := {{.LowerName}}.DeleteById(tools.String2Int(id))
	if err != nil {
		c.JSON(http.StatusOK, err)
		return
	} else {
		c.JSON(http.StatusOK, tools.GetResult(0, "Success", ""))
		return
	}
}
func Modify{{.StructName}}(c *gin.Context) {
	var {{.LowerName}} models.{{.StructName}}
	if c.ShouldBind(&{{.LowerName}}) == nil {
		fmt.Println({{.LowerName}})
		err := {{.LowerName}}.ModifyById()
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
`
const routerTemplate = `router.PUT("/{{.LowerName}}/", Add{{.StructName}})
	router.GET("/{{.LowerName}}/:id", Get{{.StructName}})
	router.GET("/{{.LowerName}}s", Get{{.StructName}})
	router.DELETE("/{{.LowerName}}/:id", Delete{{.StructName}})
	router.POST("/{{.LowerName}}/:id", Modify{{.StructName}})
	//Add router`

func GenerateCURD(structName string, primaryKey string, packageName string) (result string, err error) {
	t := template.New("curl")
	curdStrut := CurdTemplate{strings.ToLower(structName), structName, strings.ToLower(string(structName[0])), databaseIndex, primaryKey, packageName}
	//解析内容到模板
	t, err = t.Parse(curdTemplate)
	if err != nil {
		log.Fatal("Parse:", err)
		return
	}
	//将数据用到模板中
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, curdStrut); err != nil {
		log.Fatal("Execute:", err)
		return
	} else {
		result = buf.String()
	}
	return
}
func GenerateApi(structName string, packageName string) (result string, err error) {
	t := template.New("api")
	curdStrut := CurdTemplate{strings.ToLower(structName), structName, strings.ToLower(string(structName[0])), databaseIndex, "", packageName}
	//解析内容到模板
	t, err = t.Parse(apiTemplate)
	if err != nil {
		log.Fatal("Parse:", err)
		return
	}
	//将数据用到模板中
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, curdStrut); err != nil {
		log.Fatal("Execute:", err)
		return
	} else {
		result = buf.String()
	}
	return
}
func GenerateRouter(structName string, packageName string) (result string, err error) {
	t := template.New("router")
	curdStrut := CurdTemplate{strings.ToLower(structName), structName, strings.ToLower(string(structName[0])), databaseIndex, "", packageName}
	//解析内容到模板
	t, err = t.Parse(routerTemplate)
	if err != nil {
		log.Fatal("Parse:", err)
		return
	}
	//将数据用到模板中
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, curdStrut); err != nil {
		log.Fatal("Execute:", err)
		return
	} else {
		result = buf.String()
	}
	return
}
