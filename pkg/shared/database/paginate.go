package database

import "gorm.io/gorm"

func Pagination(pageData map[string]int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		page := 1
		if valPage, ok := pageData["page"]; ok && valPage > 0 {
			page = valPage
		}

		pageSize := 20
		if valPageSize, ok := pageData["per_page"]; ok && valPageSize > 0 {
			pageSize = valPageSize
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}
