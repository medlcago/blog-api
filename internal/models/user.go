package models

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           uint           `gorm:"primaryKey"`
	Username     string         `gorm:"type:string;size:50;not null;unique"`
	Email        sql.NullString `gorm:"type:string;size:256;unique;default:null"`
	Password     string         `gorm:"type:string;size:256;not null"`
	Avatar       sql.NullString `gorm:"type:string;size:256;default:null"`
	TwoFAEnabled bool           `gorm:"default:false;not null"`
	TwoFASecret  sql.NullString `gorm:"type:string;size:256;default:null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`

	Posts []Post `gorm:"foreignKey:AuthorID"`
}
