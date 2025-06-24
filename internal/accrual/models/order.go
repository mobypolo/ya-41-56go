package models

import (
	"time"
)

type Order struct {
	ID      uint    `gorm:"primaryKey"`
	Number  string  `gorm:"uniqueIndex;not null"`
	Status  string  `gorm:"type:varchar(20);not null;default:'NEW'"`
	Accrual float32 `gorm:"type:numeric(10,2);default:0"`
	Goods   []Good  `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
