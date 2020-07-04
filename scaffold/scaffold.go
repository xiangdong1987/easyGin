package scaffold

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"regexp"
	"strings"
	"xAdmin/config"
	"xAdmin/tools"
)

var mariadbHost string
var mariadbPort int
var mariadbTable string
var mariadbDatabase string
var mariadbPassword string
var mariadbUser string
var databaseIndex string
var dbConfig map[string]string

func init() {
	env := gin.Mode()
	dbConfig = config.MysqlConfigMap[env]
}

func InitDB(dbIndex string) (err error) {
	databaseIndex = dbIndex
	db := dbConfig[dbIndex]
	err = parseMysql(db)
	return
}

func parseMysql(link string) (err error) {
	mysqlRegexp := regexp.MustCompile(`(\S+):(\S+)@tcp\((\S+):(\S+)\)/(\S+)\?`)
	params := mysqlRegexp.FindStringSubmatch(link)
	for key, param := range params {
		if key == 1 {
			mariadbUser = param
		}
		if key == 2 {
			mariadbPassword = param
		}
		if key == 3 {
			mariadbHost = param
		}
		if key == 4 {
			mariadbPort = tools.String2Int(param)
		}
		if key == 5 {
			mariadbDatabase = param
		}
	}
	// Username is required
	if mariadbUser == "user" {
		err = tools.ReturnError{}.Custom(1, "Username is required! Add it with --user=name")
		return
	}
	if mariadbPassword == "" {
		err = tools.ReturnError{}.Custom(2, "Password can not be null ")
		return
	}
	if mariadbDatabase == "" {
		err = tools.ReturnError{}.Custom(2, "Database can not be null")
		return
	}
	return
}
func InitModels(table string, structName string, modelsPath string, programName string) (err error) {
	packageName := "models"
	mariadbTable = table
	columnDataTypes, err := GetColumnsFromMysqlTable(mariadbUser, mariadbPassword, mariadbHost, mariadbPort, mariadbDatabase, mariadbTable)
	if err != nil {
		fmt.Println("Error in selecting column data information from mysql information schema")
		return
	}
	// If structName is not set we need to default it
	if structName == "" {
		structName = "newstruct"
	}
	// Generate struct string based on columnDataTypes
	struc, err := Generate(*columnDataTypes, mariadbTable, structName, packageName, true, true, false)
	if err != nil {
		fmt.Println("Error in creating struct from json: " + err.Error())
		return
	}
	model := string(struc)
	//import package
	reg := regexp.MustCompile(`//packages`)
	model = reg.ReplaceAllString(model, "import (\""+programName+"/database\")")
	//get primary key
	primaryKey := getPrimaryKey(*columnDataTypes)
	curd, _ := GenerateCURD(structName, primaryKey, packageName)
	model = model + curd
	targetDirectory := modelsPath + structName + ".go"
	err = writeToFile(targetDirectory, []byte(model))
	if err != nil {
		fmt.Println(err)
		return
	}
	return
}

func InitRouter(structName string, routerPath string, packageName string) (err error) {
	routerPath = routerPath + "router.go"
	router, err := GenerateRouter(structName, packageName)
	out, isHandle, err := readFile(routerPath, "//Add router", router)
	if err != nil {
		fmt.Println("Save File fail: " + err.Error())
		return
	}
	if isHandle {
		err = writeToFile(routerPath, out)
	}
	return
}

func InitApi(structName string, apiPath string, packageName string) (err error) {
	apiPath = apiPath + strings.ToLower(structName) + ".go"
	api, err := GenerateApi(structName, packageName)
	err = writeToFile(apiPath, []byte(api))
	return
}
