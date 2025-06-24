package di

import (
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
	"ya41-56/cmd"
	"ya41-56/internal/accrual/models"
	"ya41-56/internal/accrual/worker"
	"ya41-56/internal/shared/repositories"
)

type AppContainer struct {
	OrdersRepo  repositories.Repository[models.Order]
	RewardsRepo repositories.Repository[models.RewardMechanic]

	Processor *worker.Pool
	Router    chi.Router
	Cfg       cmd.Config
	Gorm      *gorm.DB
}
