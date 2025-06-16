package bootstrap

import (
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	httpServer "net/http"
	"ya41-56/cmd"
	"ya41-56/internal/gophermart/customerror"
	dbLocal "ya41-56/internal/gophermart/db"
	"ya41-56/internal/gophermart/di"
	"ya41-56/internal/gophermart/models"
	"ya41-56/internal/gophermart/router"
	"ya41-56/internal/gophermart/services"
	"ya41-56/internal/gophermart/worker"
	"ya41-56/internal/shared/db"
	"ya41-56/internal/shared/logger"
	"ya41-56/internal/shared/repositories"
)

func Run() {
	cfg := cmd.ParseFlags()
	logger.Init(cfg.ModeLogger)

	dbConn := db.InitPostgres(&db.InitPostgresConfig{
		DSN:             cfg.DatabaseDSN,
		IsFireMigration: true,
	}, dbLocal.Migrate)

	if dbConn == nil {
		logger.L().Fatal(customerror.ErrInitDB.Error())
	}

	userRepo := repositories.NewGormRepository[models.User](dbConn)
	tokenService := services.NewTokenService(cfg.JWTSecretKey, cfg.JWTLifetime)
	orderRepo := repositories.NewGormRepository[models.Order](dbConn)

	if cfg.JWTSecretKey == "" {
		logger.L().Info(customerror.ErrEmptySecretKey.Error())
		var err error
		cfg.JWTSecretKey, err = tokenService.GenerateRandomString()
		if err != nil {
			logger.L().Fatal(customerror.ErrGenerateRandomString.Error(), zap.Error(err))
		}
	}

	fetchPool := worker.NewFetchPool(orderRepo, cfg.AccrualAddress)
	fetchPool.Start()
	defer fetchPool.Stop() //TODO давай добавим тут шатдаун из accrual

	r := router.RegisterRoutes(&di.AppContainer{
		UserRepo:  userRepo,
		OrderRepo: orderRepo,
		Auth:      services.NewAuthService(userRepo, tokenService),
		Router:    chi.NewRouter(),
		Cfg:       cfg,
		Gorm:      dbConn,
		FetchPool: fetchPool,
	})

	logger.L().Info("starting HTTP server", zap.String("addr", cfg.Address))

	err := httpServer.ListenAndServe(cfg.Address, r)
	if err != nil {
		logger.L().Fatal(customerror.ErrHTTPServer.Error(), zap.Error(err))
	}
}
