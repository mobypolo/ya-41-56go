package di

import (
	"ya41-56/cmd"
	"ya41-56/internal/gophermart/models"
	"ya41-56/internal/gophermart/services"
	"ya41-56/internal/gophermart/worker"
	"ya41-56/internal/shared/repositories"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type AppContainer struct {
	UserRepo       repositories.Repository[models.User]
	OrderRepo      repositories.Repository[models.Order]
	WithdrawalRepo repositories.Repository[models.Withdrawal]
	Auth           *services.AuthService
	Router         chi.Router
	Cfg            cmd.Config
	FetchPool      *worker.FetchPool
	Gorm           *gorm.DB
}
