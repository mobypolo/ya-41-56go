package models

import "time"

type Withdrawal struct {
	ID uint `gorm:"primaryKey"`
	// TODO: UserID take from OrderID
	UserID  uint    `gorm:"index"`
	OrderID uint    `gorm:"index"`
	Value   float32 `gorm:"default:0"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Withdrawal) TableName() string {
	return "withdrawn_gophermart"
}
