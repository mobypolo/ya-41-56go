package bootstrap

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	httpServer "net/http"
	"ya41-56/cmd"
	dbLocal "ya41-56/internal/accrual/db"
	"ya41-56/internal/accrual/di"
	"ya41-56/internal/accrual/models"
	"ya41-56/internal/accrual/router"
	"ya41-56/internal/accrual/worker"
	"ya41-56/internal/shared/db"
	"ya41-56/internal/shared/logger"
	"ya41-56/internal/shared/repositories"
	"ya41-56/internal/shared/shutdown"
)

func Run() {
	cfg := cmd.ParseFlags()
	logger.Init(cfg.ModeLogger)

	dbConn := db.InitPostgres(&db.InitPostgresConfig{
		DSN:             cfg.DatabaseDSN,
		IsFireMigration: true,
	}, dbLocal.Migrate)

	if dbConn == nil {
		logger.L().Fatal("failed to init database")
	}

	ordersRepo := repositories.NewGormRepository[models.Order](dbConn)
	rewardsRepo := repositories.NewGormRepository[models.RewardMechanic](dbConn)

	processor := worker.NewPool(ordersRepo, rewardsRepo, cfg.WorkersCount)

	r := router.RegisterRoutes(&di.AppContainer{
		OrdersRepo:  ordersRepo,
		RewardsRepo: rewardsRepo,
		Processor:   processor,
		Router:      chi.NewRouter(),
		Cfg:         cfg,
		Gorm:        dbConn,
	})

	srv := &httpServer.Server{
		Addr:    cfg.Address,
		Handler: r,
	}

	go func() {
		logger.L().Info("starting HTTP server", zap.String("addr", cfg.Address))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, httpServer.ErrServerClosed) {
			logger.L().Fatal("failed to start HTTP server", zap.Error(err))
		}
	}()

	shutdown.WaitForShutdown(srv, cfg.ShutdownTimeout,
		func() error {
			logger.L().Info("shutting down worker pool")
			processor.Stop()
			return nil
		},
	)
}
