package scaffold

import (
	"bytes"
	"log"
	"strings"
	"text/template"
)

type CurdTemplate struct {
	LowerName  string
	StructName string
	FirstChar  string
	DataIndex  string
	PrimaryKey string
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

func GenerateCURD(structName string, primaryKey string) (result string, err error) {
	t := template.New("curl")
	curdStrut := CurdTemplate{strings.ToLower(structName), structName, strings.ToLower(string(structName[0])), databaseIndex, primaryKey}
	//解析内容到模板
	t, err = t.Parse(curdTemplate)
	if err != nil {
		log.Fatal("Parse:", err)
		return
	}
	//将数据用到模板中
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, curdStrut); err != nil {
		//log.Fatal("Execute:", err)
		return
	} else {
		result = buf.String()
	}
	return
}
func GenerateApi(structName string) string {
	pk := "package handle\n"
	pks := "import (\n" +
		"	\"easyGin/models\"\n" +
		"	\"easyGin/tools\"\n" +
		"	\"fmt\"\n" +
		"	\"github.com/gin-gonic/gin\"\n" +
		"	\"net/http\"\n" +
		")\n"
	add := "func Add" + structName + "(c *gin.Context) {\n" +
		"	var " + strings.ToLower(structName) + " models." + structName + "\n" +
		"	if c.ShouldBind(&" + strings.ToLower(structName) + ") == nil {\n" +
		"		fmt.Println(" + strings.ToLower(structName) + ")\n" +
		"		err := " + strings.ToLower(structName) + ".Insert()\n" +
		"		if err != nil {\n" +
		"			c.JSON(http.StatusOK, err)\n" +
		"			return		\n" +
		"		} else {\n" +
		"			c.JSON(http.StatusOK, tools.GetResult(0, \"Success\", \"\"))\n" +
		"			return\n" +
		"		}\n" +
		"	}\n" +
		"	c.JSON(http.StatusOK, tools.GetResult(0, \"Params is wrong\", \"\"))\n" +
		"	return\n" +
		"}\n"
	get := "func Get" + structName + "(c *gin.Context) {\n" +
		"	id := c.Param(\"id\")\n" +
		"	var " + strings.ToLower(structName) + " models." + structName + "\n" +
		"	if id != \"\" {" +
		"		err := " + strings.ToLower(structName) + ".GetById(tools.String2Int(id))\n" +
		"		if err != nil {\n" +
		"			c.JSON(http.StatusOK, err)\n" +
		"			return\n" +
		"		} else {\n" +
		"			c.JSON(http.StatusOK, tools.GetResult(0, \"Success\", gin.H{\n" +
		"				\"" + strings.ToLower(structName) + "\": " + strings.ToLower(structName) + ",\n" +
		"			}))\n" +
		"			return\n" +
		"		}\n" +
		"	} else {\n" +
		"		page := c.DefaultQuery(\"page\", \"1\")\n" +
		"		pageSize := c.DefaultQuery(\"page_size\", \"10\")\n" +
		"		list, pageInfo, err := " + strings.ToLower(structName) + ".GetList(tools.String2Int(page), tools.String2Int(pageSize))\n" +
		"		if err != nil {\n" +
		"			c.JSON(http.StatusOK, err)\n" +
		"			return\n" +
		"		} else {\n" +
		"			c.JSON(http.StatusOK, tools.GetResult(0, \"Success\", gin.H{\n" +
		"				\"list\":      list,\n" +
		"				\"page_info\": pageInfo,\n" +
		"			}))\n" +
		"			return\n" +
		"		}\n" +
		"	}\n" +
		"}\n"
	delete := "func Delete" + structName + "(c *gin.Context) {\n" +
		"	id := c.Param(\"id\")\n" +
		"	var " + strings.ToLower(structName) + " models." + structName + "\n" +
		"	err := " + strings.ToLower(structName) + ".DeleteById(tools.String2Int(id))\n" +
		"	if err != nil {\n" +
		"		c.JSON(http.StatusOK, err)\n" +
		"		return\n" +
		"	} else {\n" +
		"		c.JSON(http.StatusOK, tools.GetResult(0, \"Success\", \"\"))\n" +
		"		return\n" +
		"	}\n" +
		"}\n"
	update := "func Modify" + structName + "(c *gin.Context) {\n" +
		"	var person models." + structName + "\n" +
		"	if c.ShouldBind(&" + strings.ToLower(structName) + ") == nil {\n" +
		"		fmt.Println(" + strings.ToLower(structName) + ")\n" +
		"		err := " + strings.ToLower(structName) + ".ModifyById()\n" +
		"		if err != nil {\n" +
		"			c.JSON(http.StatusOK, err)\n" +
		"			return\n" +
		"		} else {\n" +
		"			c.JSON(http.StatusOK, tools.GetResult(0, \"Success\", \"\"))\n" +
		"			return\n" +
		"		}\n" +
		"	}\n" +
		"	c.JSON(http.StatusOK, tools.GetResult(0, \"Params is wrong\", \"\"))\n" +
		"	return\n" +
		"}\n"
	return pk + pks + add + get + delete + update
}
func GenerateRouter(structName string) string {
	router := "router.POST(\"/" + strings.ToLower(string(structName)) + "/\", Add" + structName + ")\n" +
		"	router.GET(\"/" + strings.ToLower(string(structName)) + "/:id\", Get" + structName + ")\n" +
		"	router.GET(\"/" + strings.ToLower(string(structName)) + "s\", Get" + structName + ")\n" +
		"	router.DELETE(\"/" + strings.ToLower(string(structName)) + "/:id\", Delete" + structName + ")\n" +
		"	router.PUT(\"/" + strings.ToLower(string(structName)) + "/:id\", Modify" + structName + ")\n" +
		"	//Add router\n"
	return router
}
