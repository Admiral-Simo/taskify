package models

import "time"

const (
	LowPriority    = "L"
	MediumPriority = "M"
	HighPriority   = "H"
)

type Task struct {
	ID        int64     `gorm:"primaryKey"`
	Title     string    `gorm:"not null"`
	Done      bool      `gorm:"default:false"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	Priority  string    `gorm:"default:L"`
}
