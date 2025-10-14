package models

import "time"

type User struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	Email     string    `gorm:"unique" json:"email"`
	Password  string    `gorm:"not null" json:"-"`
	Verified  bool      `gorm:"not null" json:"-"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}
