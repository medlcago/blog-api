package models

import (
	"time"
)

type ReactionType string

const (
	ReactionLike    ReactionType = "like"
	ReactionDislike ReactionType = "dislike"
	ReactionLove    ReactionType = "love"
	ReactionLaugh   ReactionType = "laugh"
	ReactionSad     ReactionType = "sad"
	ReactionAngry   ReactionType = "angry"
	ReactionFire    ReactionType = "fire"
)

var (
	AllowedReactions = []ReactionType{
		ReactionLike, ReactionDislike, ReactionLove,
		ReactionLaugh, ReactionSad, ReactionAngry, ReactionFire,
	}
)

type Reaction struct {
	ID     uint         `gorm:"primaryKey"`
	UserID uint         `gorm:"not null;uniqueIndex:idx_user_target"`
	Type   ReactionType `gorm:"type:varchar(20);not null"`

	TargetID   uint   `gorm:"not null;uniqueIndex:idx_user_target"`
	TargetType string `gorm:"size:50;not null;uniqueIndex:idx_user_target"`

	CreatedAt time.Time
	User      User `gorm:"foreignKey:UserID"`
}
