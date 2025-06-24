package models

import (
	"time"
	"ya41-56/internal/shared/models"
)

type RewardMechanic struct {
	ID         uint              `gorm:"primaryKey"`
	Match      string            `gorm:"type:text;not null"`
	Reward     float32           `gorm:"not null"`
	RewardType models.RewardType `gorm:"type:varchar(10);not null"` // 'pt', '%'

	CreatedAt time.Time
	UpdatedAt time.Time
}
