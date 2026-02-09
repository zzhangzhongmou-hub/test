package models

import "time"

type Homework struct {
	ID            uint      `gorm:"primarykey" json:"id"`
	Title         string    `gorm:"size:200;not null" json:"title"`
	Description   string    `gorm:"type:text;not null" json:"description"`
	Department    string    `gorm:"size:20;not null;index" json:"department"`
	CreatorID     uint      `gorm:"not null;" json:"creator_id"`
	Deadline      time.Time `gorm:"not null" json:"deadline"`
	AllowLate     bool      `gorm:"default:false" json:"allow_late"`
	Version       int       `gorm:"default:1" json:"version"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Type          string    `gorm:"default:'homework';index"`
	ReviewerCount int       `gorm:"default:1"`
}

func (h *Homework) IsExam() bool {
	return h.Type == "exam"
}
