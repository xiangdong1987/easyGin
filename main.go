package main

import (
	"easyGin/api"
	"easyGin/api/service"
	"easyGin/config"
	"easyGin/database"
	"easyGin/scaffold"
	"fmt"
	middleware "github.com/deepmap/oapi-codegen/pkg/gin-middleware"
	"github.com/droundy/goopt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
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

func NewGinServer(demo *api.Demo, port int) *http.Server {
	swagger, err := service.GetSwagger()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading swagger spec\n: %s", err)
		os.Exit(1)
	}

	// Clear out the servers array in the swagger spec, that skips validating
	// that server names match. We don't know how this thing will be run.
	swagger.Servers = nil

	// This is how you set up a basic chi router
	r := gin.Default()
	// Use our validation middleware to check all requests against the
	// OpenAPI schema.
	r.Use(middleware.OapiRequestValidator(swagger))
	// We now register our petStore above as the handler for the interface
	r = service.RegisterHandlers(r, demo)

	s := &http.Server{
		Handler: r,
		Addr:    fmt.Sprintf("0.0.0.0:%d", port),
	}
	return s
}

func InitServer(port int) {
	// Create an instance of our handler which satisfies the generated interface
	demo := api.NewDemo()
	s := NewGinServer(demo, port)
	// And we serve HTTP until the world ends.
	log.Fatal(s.ListenAndServe())
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
		InitServer(8082)
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
}
