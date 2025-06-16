package models

import "time"

const (
	OrderStatusNew = "NEW"
	//OrderStatusProcessing  = "PROCESSING"
	//OrderStatusProcessed   = "PROCESSED"
	OrderStatusFailedFetch = "FAILED_FETCH"
)

type Order struct {
	ID      uint    `gorm:"primaryKey"`
	UserID  uint    `gorm:"index"`
	Number  string  `gorm:"uniqueIndex;not null"`
	Status  string  `gorm:"default:'NEW'"`
	Accrual float32 `gorm:"default:0"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Order) TableName() string {
	return "orders_gophermart"
}
