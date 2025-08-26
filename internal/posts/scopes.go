package posts

import (
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func OrderScope(orderBy, sort string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if orderBy == "" {
			orderBy = "created_at"
		}

		desc := strings.ToLower(sort) == "desc"
		return db.Order(clause.OrderByColumn{Column: clause.Column{Name: orderBy}, Desc: desc})
	}
}

func PaginationScope(limit, offset int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		switch {
		case limit > 100:
			limit = 100
		case limit <= 0:
			limit = 30
		}

		db.Limit(limit)

		if offset > 0 {
			db = db.Offset(offset)
		}
		return db
	}
}
