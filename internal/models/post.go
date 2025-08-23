package models

import (
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	Title    string `gorm:"size:255;not null"`
	Content  string `gorm:"type:text;not null"`
	AuthorID uint   `gorm:"index;not null"`

	Author   User         `gorm:"foreignKey:AuthorID"`
	Entities []PostEntity `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE"`
}
