package di

import (
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
	"ya41-56/cmd"
	"ya41-56/internal/gophermart/models"
	"ya41-56/internal/gophermart/services"
	"ya41-56/internal/shared/repositories"
)

type AppContainer struct {
	UserRepo repositories.Repository[models.User]
	Auth     *services.AuthService
	Router   chi.Router
	Cfg      cmd.Config
	Gorm     *gorm.DB
}
