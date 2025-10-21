package models

import "time"

type User struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	Email     string    `gorm:"unique" json:"email"`
	Password  string    `gorm:"not null" json:"-"`
	Verified  bool      `gorm:"not null" json:"verified"`
	Profile   string    `json:"profile"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}
