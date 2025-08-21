package database

import "gorm.io/gorm"

type Filter interface {
	Apply(query *gorm.DB) *gorm.DB
}
