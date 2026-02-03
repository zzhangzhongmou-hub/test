package models

import "time"

type Exam struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	Title       string    `gorm:"size:200;not null" json:"title"`
	Description string    `gorm:"type:text;not null" json:"description"`
	Department  string    `gorm:"size:20;not null;index" json:"department"`
	CreatorID   uint      `gorm:"not null" json:"creator_id"`
	Deadline    time.Time `gorm:"not null" json:"deadline"`
	Duration    int       `json:"duration"`
	TotalScore  int       `gorm:"default:100" json:"total_score"`
	Status      string    `gorm:"size:20;default:'pending'" json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
