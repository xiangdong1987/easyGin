package main

import (
	"easyGin/config"
	orm "easyGin/database"
	"easyGin/scaffold"
	"github.com/droundy/goopt"
	"go-pc_home/router"
)

var ifScaffold = goopt.Int([]string{"--ifScaffold"}, 0, "if use scaffold")
var databaseIndex = goopt.String([]string{"--index"}, "", "which config is used")
var table = goopt.String([]string{"-t", "--table"}, "", "Table to build struct from")
var structName = goopt.String([]string{"--struct"}, "", "name to set for struct")

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
			for _, value := range orm.DBList {
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
	if *databaseIndex == "" {
		panic("database index can not null")
	}
	if *table == "" {
		panic("table can not null")
	}
	if *structName == "" {
		panic("structName can not null")
	}
	scaffold.InitDB(*databaseIndex)
	scaffold.InitModels(*table, *structName)
	scaffold.InitApi(*structName)
	scaffold.InitRouter(*structName)
}
