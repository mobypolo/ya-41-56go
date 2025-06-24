package models

import (
	"time"
)

type UserStatus int

const (
	UserStatusDisabled UserStatus = 0
	UserStatusActive   UserStatus = 1
)

type User struct {
	ID           uint       `gorm:"primaryKey"`
	Login        string     `gorm:"uniqueIndex;not null"`
	Password     string     `gorm:"-" json:"-"`
	PasswordHash string     `gorm:"not null" json:"-"`
	Status       UserStatus `gorm:"default:0"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
