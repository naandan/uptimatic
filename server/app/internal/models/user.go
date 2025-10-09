package models

import "time"

type User struct {
	ID        uint      `gorm:"primary_key"`
	Email     string    `gorm:"unique"`
	Password  string    `gorm:"not null"`
	Verified  bool      `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
