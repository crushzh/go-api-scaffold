package model

import "time"

// User is the user model (JWT authentication)
type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username" gorm:"size:50;uniqueIndex;not null"`
	Password  string    `json:"-" gorm:"size:100;not null"` // Hidden from JSON output
	Role      string    `json:"role" gorm:"size:20;default:user"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}
