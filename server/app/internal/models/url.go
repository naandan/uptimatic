package models

import (
	"time"

	"github.com/google/uuid"
)

type URL struct {
	ID          uint       `gorm:"primary_key"`
	PublicID    uuid.UUID  `gorm:"not null;unique"`
	UserID      uint       `gorm:"not null"`
	Label       string     `gorm:"not null"`
	URL         string     `gorm:"not null"`
	Interval    int        `gorm:"not null"`
	Active      bool       `gorm:"not null"`
	LastChecked *time.Time `gorm:"null"`
	CreatedAt   time.Time  `gorm:"autoCreateTime"`

	User User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type StatusLog struct {
	ID           uint      `gorm:"primary_key"`
	URLID        uint      `gorm:"not null"`
	Status       string    `gorm:"not null"`
	ResponseTime int64     `gorm:"not null"`
	CheckedAt    time.Time `gorm:"autoCreateTime"`
}

type UptimeStat struct {
	BucketStart   time.Time `json:"bucket_start"`
	TotalChecks   int       `json:"total_checks"`
	UpChecks      int       `json:"up_checks"`
	UptimePercent float64   `json:"uptime_percent"`
}
