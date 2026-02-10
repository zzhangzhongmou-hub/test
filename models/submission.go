package models

import (
	"time"

	"gorm.io/gorm"
)

type Submission struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	Type        string         `gorm:"size:20;not null;index" json:"type"` //考核与作业
	HomeworkID  uint           `gorm:"index" json:"homework_id"`
	ExamID      uint           `gorm:"index" json:"exam_id"`
	StudentID   uint           `gorm:"not null;index" json:"student_id"`
	Content     string         `gorm:"type:text;not null" json:"content"`
	IsLate      bool           `gorm:"default:false" json:"is_late"`
	FileURL     string         `gorm:"size:500" json:"file_url"`
	Score       *int           `json:"score"`
	Comment     string         `gorm:"type:text" json:"comment"`
	IsExcellent bool           `gorm:"default:false;index" json:"is_excellent"`
	ReviewerID  *uint          `json:"reviewer_id"`
	SubmittedAt time.Time      `gorm:"not null" json:"submitted_at"`
	ReviewedAt  *time.Time     `json:"reviewed_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Status      string         `json:"status"`
	Homework    Homework       `gorm:"foreignKey:HomeworkID" json:"homework,omitempty"`
	Student     User           `gorm:"foreignKey:StudentID" json:"student,omitempty"`
}
