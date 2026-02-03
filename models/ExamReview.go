package models

import "time"

type ExamReview struct {
	ID         uint       `gorm:"primarykey" json:"id"`
	ExamID     uint       `gorm:"not null;index" json:"exam_id"`
	StudentID  uint       `gorm:"not null;index" json:"student_id"`
	ReviewerID uint       `gorm:"not null;index" json:"reviewer_id"`
	Score      int        `gorm:"not null" json:"score"`
	Comment    string     `gorm:"type:text" json:"comment"`
	Status     string     `gorm:"size:20;default:'pending'" json:"status"`
	AssignedAt time.Time  `gorm:"not null" json:"assigned_at"`
	ReviewedAt *time.Time `json:"reviewed_at"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}