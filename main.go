package main

import (
	"easyGin/config"
	"easyGin/database"
	"easyGin/router"
	"easyGin/scaffold"
	"github.com/droundy/goopt"
)

var ifScaffold = goopt.Int([]string{"--ifScaffold"}, 0, "if use scaffold")
var table = goopt.String([]string{"-t", "--table"}, "", "Table to build struct from")
var structName = goopt.String([]string{"--struct"}, "", "name to set for struct")
var program = goopt.String([]string{"--program"}, "", "name to set for program")
var dbIndex = goopt.String([]string{"--dbIndex"}, "default", "name to set for program")

func init() {
	// Setup goopts
	goopt.Description = func() string {
		return "EasyGin init restful api"
	}
	goopt.Version = "0.0.2"
	goopt.Summary = "main --struct structName --index db_index databaseName --table tableName"
	//Parse options
	goopt.Parse(nil)

}

func main() {
	if *ifScaffold == 0 {
		defer func() {
			for _, value := range database.DBList {
				value.Close()
			}
		}()
		defer config.RedisObj.Close()
		config.SetRedisObj()
		r := router.InitRouter()
		r.Run(":8082")
	} else {
		scaffoldInit()
	}
}

func scaffoldInit() {
	if *table == "" {
		panic("table can not null")
	}
	if *structName == "" {
		panic("structName can not null")
	}
	if *dbIndex == "" {
		panic("dbIndex can not null")
	}
	if *program == "" {
		*program = "easyGin"
	}
	scaffold.InitDB(*dbIndex)
	scaffold.InitModels(*table, *structName, config.ModelPath, *program)
	scaffold.InitApi(*structName, config.ApiPath, *program)
	scaffold.InitRouter(*structName, config.RouterPath, *program)
}