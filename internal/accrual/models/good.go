package models

type Good struct {
	ID          uint    `gorm:"primaryKey"`
	OrderID     uint    `gorm:"index"`
	Description string  `gorm:"not null"`
	Price       float64 `gorm:"not null"`
}
