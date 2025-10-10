package models

import "time"

type URL struct {
	ID          uint       `gorm:"primary_key"`
	UserID      uint       `gorm:"not null"`
	Label       string     `gorm:"not null"`
	URL         string     `gorm:"not null"`
	Interval    int        `gorm:"not null"`
	Active      bool       `gorm:"not null"`
	LastChecked *time.Time `gorm:"null"`
	CreatedAt   time.Time  `gorm:"autoCreateTime"`
}

type StatusLog struct {
	ID           uint      `gorm:"primary_key"`
	URLID        uint      `gorm:"not null"`
	Status       string    `gorm:"not null"`
	ResponseTime int64     `gorm:"not null"`
	CheckedAt    time.Time `gorm:"autoCreateTime"`
}
