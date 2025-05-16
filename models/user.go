package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username     string    `gorm:"uniqueIndex;not null" json:"username"`
	Password     string    `gorm:"not null" json:"-"`
	LastLogin    time.Time `json:"last_login"`
	LastReward   time.Time `json:"last_reward"`
	RewardStreak int       `gorm:"default:0" json:"reward_streak"`
	Coins        int       `gorm:"default:0" json:"coins"`
}
