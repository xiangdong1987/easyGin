package scaffold

import (
	"easyGin/config"
	"easyGin/tools"
	"fmt"
	"github.com/Shelnutt2/db2struct"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"regexp"
	"strings"
)

var mariadbHost string
var mariadbPort int
var mariadbTable string
var mariadbDatabase string
var mariadbPassword string
var mariadbUser string

var databaseIndex = "test"

func init() {
	env := gin.Mode()
	dbConfig := config.MysqlConfigMap[env]
	db := dbConfig[databaseIndex]
	fmt.Println(db)
	parseMysql(db)
}
func parseMysql(link string) {
	flysnowRegexp := regexp.MustCompile(`(\S+):(\S+)@tcp\((\S+):(\S+)\)/(\S+)\?`)
	params := flysnowRegexp.FindStringSubmatch(link)
	fmt.Println(params)
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
}
func InitModels(table string, structName string, packageName string, jsonAnnotation bool, gormAnnotation bool, gureguTypes bool) {
	mariadbTable = table
	// Username is required
	if mariadbUser == "user" {
		fmt.Println("Username is required! Add it with --user=name")
		return
	}
	if mariadbPassword == "" {
		fmt.Println("Password can not be null ")

	}
	if mariadbDatabase == "" {
		fmt.Println("Database can not be null")
		return
	}
	if mariadbTable == "" {
		fmt.Println("Table can not be null")
		return
	}
	columnDataTypes, err := db2struct.GetColumnsFromMysqlTable(mariadbUser, mariadbPassword, mariadbHost, mariadbPort, mariadbDatabase, mariadbTable)
	if err != nil {
		fmt.Println("Error in selecting column data information from mysql information schema")
		return
	}
	// If structName is not set we need to default it
	if structName == "" {
		structName = "newstruct"
	}
	// If packageName is not set we need to default it
	if packageName == "" {
		packageName = "newpackage"
	}
	// Generate struct string based on columnDataTypes
	struc, err := db2struct.Generate(*columnDataTypes, mariadbTable, structName, packageName, jsonAnnotation, gormAnnotation, gureguTypes)
	if err != nil {
		fmt.Println("Error in creating struct from json: " + err.Error())
		return
	}
	model := string(struc)
	model = model + GenerateCURD(structName, "id")
	targetDirectory := config.ModelPath + structName + ".go"
	file, err := os.OpenFile(targetDirectory, os.O_WRONLY|os.O_CREATE, os.ModeAppend)
	if err != nil {
		fmt.Println("Open File fail: " + err.Error())
		return
	}
	length, err := file.WriteString(model)
	if err != nil {
		fmt.Println("Save File fail: " + err.Error())
		return
	}
	fmt.Printf("wrote %d bytes\n", length)
	//fmt.Println(string(struc))
}

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
