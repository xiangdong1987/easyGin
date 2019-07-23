package scaffold

import (
	"fmt"
	"github.com/Shelnutt2/db2struct"
	"github.com/gin-gonic/gin"
	"go-pc_home/config"

	_ "github.com/go-sql-driver/mysql"
	"github.com/howeyc/gopass"
)

var mariadbHost string
var mariadbHostPassed string
var mariadbPort int
var mariadbTable string
var mariadbDatabase string
var mariadbPassword *string
var mariadbUser string

func init() {
	env := gin.Mode()
	dbConfig := config.MysqlConfigMap[env]
	db := dbConfig["test"]
	
}

func InitModels(structName string, packageName string, jsonAnnotation bool, gormAnnotation bool, gureguTypes bool) {

	// Username is required
	if mariadbUser == "user" {
		fmt.Println("Username is required! Add it with --user=name")
		return
	}
	if mariadbPassword != nil && *mariadbPassword == "" {
		fmt.Print("Password: ")
		pass, err := gopass.GetPasswd()
		stringPass := string(pass)
		mariadbPassword = &stringPass
		if err != nil {
			fmt.Println("Error reading password: " + err.Error())
			return
		}
	}
	if mariadbDatabase == "" {
		fmt.Println("Database can not be null")
		return
	}
	if mariadbTable == "" {
		fmt.Println("Table can not be null")
		return
	}
	columnDataTypes, err := db2struct.GetColumnsFromMysqlTable(mariadbUser, *mariadbPassword, mariadbHost, mariadbPort, mariadbDatabase, mariadbTable)
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
	fmt.Println(string(struc))
}
