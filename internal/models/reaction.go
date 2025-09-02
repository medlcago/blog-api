package models

import "time"

type Reaction struct {
	ID     uint `gorm:"primaryKey"`
	UserID uint `gorm:"not null;uniqueIndex:idx_user_target"`

	TargetID       uint   `gorm:"not null;uniqueIndex:idx_user_target"`
	TargetType     string `gorm:"size:50;not null;uniqueIndex:idx_user_target"`
	ReactionTypeID uint   `gorm:"not null"`

	CreatedAt time.Time
	UpdatedAt time.Time

	User         User         `gorm:"foreignKey:UserID"`
	ReactionType ReactionType `gorm:"foreignKey:ReactionTypeID"`
}

type UserReaction struct {
	TargetID uint   `json:"-"`
	Type     string `json:"type"`
	Icon     string `json:"icon"`
	IsActive bool   `json:"is_active"`
}

type ReactionStat struct {
	TargetID uint   `json:"-"`
	Type     string `json:"type"`
	Count    int64  `json:"count"`
	Icon     string `json:"icon"`
	IsActive bool   `json:"is_active"`
}
