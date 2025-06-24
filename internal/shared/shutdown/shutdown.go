package shutdown

import (
	"context"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"ya41-56/internal/shared/logger"
)

func WaitForShutdown(server *http.Server, ShutdownTimeout time.Duration, cleanupFns ...func() error) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	logger.L().Info("shutting down HTTP server...")

	ctx, cancel := context.WithTimeout(context.Background(), ShutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.L().Error("graceful shutdown failed", zap.Error(err))
	} else {
		logger.L().Info("server stopped cleanly")
	}

	for _, fn := range cleanupFns {
		if err := fn(); err != nil {
			logger.L().Error("cleanup function failed", zap.Error(err))
		}
	}
}
