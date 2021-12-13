package api

import (
	"easyGin/api/service"
	"easyGin/handle"
	"github.com/gin-gonic/gin"
)

//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen --config=types.cfg.yaml person.ymal
//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen --config=server.cfg.yaml person.ymal

type Demo struct {
}

func NewDemo() *Demo {
	return &Demo{}
}

func (d *Demo) FindPerson(c *gin.Context, params service.FindPersonParams) {
	return
}

func (d *Demo) AddPerson(c *gin.Context) {
	handle.AddPerson(c)
	return
}

func (d *Demo) DeletePerson(c *gin.Context, id int64) {
	handle.DeletePerson(c, id)
	return
}

func (d *Demo) FindPersonByID(c *gin.Context, id int64) {
	handle.GetPerson(c, id)
}
