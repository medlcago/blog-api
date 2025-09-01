package models

import "time"

type ReactionType struct {
	ID       uint   `gorm:"primaryKey"`
	Name     string `gorm:"unique;size:50;not null"`
	Icon     string `gorm:"size:50;not null"`
	IsActive bool   `gorm:"not null;default:true"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
