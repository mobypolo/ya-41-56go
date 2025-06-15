package bootstrap

import (
	httpServer "net/http"
	"ya41-56/cmd"
	dbLocal "ya41-56/internal/gophermart/db"
	"ya41-56/internal/gophermart/di"
	commonErrors "ya41-56/internal/gophermart/errors"
	"ya41-56/internal/gophermart/models"
	"ya41-56/internal/gophermart/router"
	"ya41-56/internal/gophermart/services"
	"ya41-56/internal/shared/db"
	"ya41-56/internal/shared/logger"
	"ya41-56/internal/shared/repositories"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func Run() {
	cfg := cmd.ParseFlags()
	logger.Init(cfg.ModeLogger)

	dbConn := db.InitPostgres(&db.InitPostgresConfig{
		DSN:             cfg.DatabaseDSN,
		IsFireMigration: true,
	}, dbLocal.Migrate)

	if dbConn == nil {
		logger.L().Fatal(commonErrors.ErrInitDB.Error())
	}

	userRepo := repositories.NewGormRepository[models.User](dbConn)

	if cfg.JWTSecretKey == "" {
		logger.L().Fatal(commonErrors.ErrEmptySecretKey.Error())
	}

	r := router.RegisterRoutes(&di.AppContainer{
		UserRepo: userRepo,
		Auth:     services.NewAuthService(userRepo, services.NewTokenService(cfg.JWTSecretKey, cfg.JWTLifetime)),
		Router:   chi.NewRouter(),
		Cfg:      cfg,
		Gorm:     dbConn,
	})

	logger.L().Info("starting HTTP server", zap.String("addr", cfg.Address))

	err := httpServer.ListenAndServe(cfg.Address, r)
	if err != nil {
		logger.L().Fatal(commonErrors.ErrHTTPServer.Error(), zap.Error(err))
	}
}
