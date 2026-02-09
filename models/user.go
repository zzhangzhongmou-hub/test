package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	Username   string         `gorm:"uniqueIndex;size:50;not null" json:"firstname"`
	Nickname   string         `gorm:"size:50;not null" json:"nickname"`
	Password   string         `gorm:"size:255;not null"`
	Role       string         `gorm:"size:20;not null" json:"role"`
	Department string         `gorm:"size:20;not null" json:"department"`
	Email      string         `gorm:"size:100" json:"email"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}
