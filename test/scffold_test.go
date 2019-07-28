package test

import (
	"easyGin/scaffold"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestGenerateCURD(t *testing.T) {
	Convey("生成Curd", t, func() {
		scaffold.InitDB("test")
		Convey("正确生成curd", func() {
			x, _ := scaffold.GenerateCURD("Person", "id")
			So(x, ShouldEqual, `func (p *Person) Insert() (err error) {
	dbe, err := database.Database("test") 
	if err != nil { 
		return 
	} 
	err = dbe.Create(p).Error 
	return 
}
func (p *Person) GetById(id int) (err error) {
	dbe, err := database.Database("test")
	if err != nil {
		return
	}
	err = dbe.Where("id=?", id).First(p).Error
	return
}
func (p *Person) GetList(currentPage int, pageSize int) (list []Person, pageInfo database.PageInfo, err error) { 
	dbe, err := database.Database("test")
	if err != nil {
		return
	}
	start := (currentPage - 1) * pageSize
	dbList := dbe.Limit(pageSize).Offset(start).Order("id DESC")
	err = dbList.Find(&list).Error
	pageInfo = pageInfo.GetPageInfo(dbe, &Person{}, currentPage, pageSize)
	return
}
func (p *Person) DeleteById(id int) (err error) {
	dbe, err := database.Database("test")
	if err != nil {
		return
	}
	err = dbe.Where("id=?", id).Delete(p).Error
	return
}
func (p *Person) ModifyById() (err error) {
	dbe, err := database.Database("test")
	if err != nil {
		return
	}
	err = dbe.Save(p).Error
	return
}`)
		})
	})
}
