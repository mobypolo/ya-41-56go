package db

import (
	"ya41-56/internal/gophermart/models"
	"ya41-56/internal/shared/logger"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(
		models.User{},
		models.Order{},
		models.Withdrawal{},
	); err != nil {
		logger.L().Info("auto migration to postgres failed")
		return err
	}

	return nil
}
