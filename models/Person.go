package models

import "easyGin/database"

type Person struct {
	Age   int    `gorm:"column:age" json:"age" form:"age"`
	Ctime int    `gorm:"column:ctime" json:"ctime" form:"ctime"`
	ID    int    `gorm:"column:id;primary_key" json:"id;primary_key" form:"id"`
	Name  string `gorm:"column:name" json:"name" form:"name"`
	Sex   int    `gorm:"column:sex" json:"sex" form:"sex"`
}

// TableName sets the insert table name for this struct type
func (p *Person) TableName() string {
	return "person"
}
func (p *Person) Insert() (err error) {
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
}
