package models

import (
	"database/sql"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string         `gorm:"type:string;size:50;not null;unique"`
	Email    sql.NullString `gorm:"type:string;size:256;unique;default:null"`
	Password string         `gorm:"type:string;size:256;not null"`

	Posts []Post `gorm:"foreignKey:AuthorID"`
}
