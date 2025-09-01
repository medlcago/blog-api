package models

import "gorm.io/gorm"

type PostEntity struct {
	gorm.Model
	PostID uint    `gorm:"index;not null"`
	Offset int     `gorm:"not null"`
	Length int     `gorm:"not null"`
	Type   string  `gorm:"size:50;not null"`
	URL    *string `gorm:"size:500"`

	Post Post `gorm:"foreignKey:PostID"`
}

func (e *PostEntity) TableName() string {
	return "post_entities"
}
