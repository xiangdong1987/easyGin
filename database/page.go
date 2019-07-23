package database

import "github.com/jinzhu/gorm"

type PageInfo struct {
	CurrentPage int `json:"page" form:"page"`
	TotalNum    int `json:"total_num" form:"total_num"`
	PageSize    int `json:"page_size" form:"page_size"`
}

func (p PageInfo) GetPageInfo(db *gorm.DB, table interface{}, currentPage int, pageSize int) PageInfo {
	p.CurrentPage = currentPage
	p.PageSize = pageSize
	//获取总数
	db.Model(table).Count(&p.TotalNum)
	return p
}
