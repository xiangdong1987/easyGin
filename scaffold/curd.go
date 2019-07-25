package scaffold

import (
	"strings"
)

func GenerateCURD(structName string, primaryKey string) string {
	insert := "func (" + strings.ToLower(string(structName[0])) + " *" + structName + ") Insert() (err error) {\n" +
		"	dbe, err := database.Database(\"" + databaseIndex + "\") \n" +
		"	if err != nil { \n" +
		"		return \n" +
		"	} \n" +
		"	err = dbe.Create(" + strings.ToLower(string(structName[0])) + ").Error \n" +
		"	return \n" +
		"}\n"
	getOne := "func (" + strings.ToLower(string(structName[0])) + " *" + structName + ") GetById(id int) (err error) {\n" +
		"	dbe, err := database.Database(\"" + databaseIndex + "\")\n" +
		"	if err != nil {\n" +
		"		return\n" +
		"	}\n" +
		"	err = dbe.Where(\"" + primaryKey + "=?\", id).First(" + strings.ToLower(string(structName[0])) + ").Error\n" +
		"	return\n" +
		"}\n"
	getList := "func (" + strings.ToLower(string(structName[0])) + " *" + structName + ") GetList(currentPage int, pageSize int) (list []" + structName + ", pageInfo database.PageInfo, err error) { \n" +
		"	dbe, err := database.Database(\"" + databaseIndex + "\")\n" +
		"	if err != nil {\n" +
		"		return\n" +
		"	}\n" +
		"	start := (currentPage - 1) * pageSize\n" +
		"	dbList := dbe.Limit(pageSize).Offset(start).Order(\"" + primaryKey + " DESC\")\n" +
		"	err = dbList.Find(&list).Error\n" +
		"	pageInfo = pageInfo.GetPageInfo(dbe, &" + structName + "{}, currentPage, pageSize)\n" +
		"	return\n" +
		"}\n"
	deleteID := "func (" + strings.ToLower(string(structName[0])) + " *" + structName + ") DeleteById(id int) (err error) {\n" +
		"	dbe, err := database.Database(\"" + databaseIndex + "\")\n" +
		"	if err != nil {\n" +
		"		return\n" +
		"	}\n" +
		"	err = dbe.Where(\"" + primaryKey + "=?\", id).Delete(" + strings.ToLower(string(structName[0])) + ").Error\n" +
		"	return\n" +
		"}\n"
	update := "func (" + strings.ToLower(string(structName[0])) + " *" + structName + ") ModifyById() (err error) {\n" +
		"	dbe, err := database.Database(\"" + databaseIndex + "\")\n" +
		"	if err != nil {\n" +
		"		return\n" +
		"	}\n" +
		"	err = dbe.Save(" + strings.ToLower(string(structName[0])) + ").Error\n" +
		"	return\n" +
		"}\n"
	//println(insert, getOne, getList, delete, update)
	return insert + getOne + getList + deleteID + update
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
