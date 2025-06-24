package router

import (
	"net/http"
	"ya41-56/internal/accrual/di"
	"ya41-56/internal/accrual/handlers"
	sharedHandlers "ya41-56/internal/shared/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func RegisterRoutes(appContainer *di.AppContainer) http.Handler {
	ordersHandler := handlers.NewOrdersHandler(appContainer.OrdersRepo, appContainer.Processor)
	rewardsHandler := handlers.NewRewardsHandler(appContainer.RewardsRepo)

	appContainer.Router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   appContainer.Cfg.CorsOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	appContainer.Router.Get("/ping", sharedHandlers.PingHandler(appContainer.Gorm))

	appContainer.Router.Route("/api", func(r chi.Router) {
		r.Get("/orders/{number}", ordersHandler.GetOrderByNumber)
		r.Post("/orders", ordersHandler.CreateOrder)
		r.Post("/goods", rewardsHandler.CreateRewards)
	})

	return appContainer.Router
}
