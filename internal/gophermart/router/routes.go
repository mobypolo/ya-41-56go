package router

import (
	"net/http"
	"ya41-56/internal/gophermart/di"
	"ya41-56/internal/gophermart/handlers"
	sharedHandlers "ya41-56/internal/shared/handlers"
	"ya41-56/internal/shared/http/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func RegisterRoutes(appContainer *di.AppContainer) http.Handler {
	authMiddleware := middleware.New(appContainer.Auth, appContainer.UserRepo)

	authHandler := handlers.NewAuthHandler(appContainer.Auth)
	usersHandler := handlers.NewUsersHandler(appContainer.Auth)

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
		r.Route("/user", func(r chi.Router) {
			r.Post("/login", authHandler.Login)
			r.Post("/register", authHandler.Register)

			r.Group(func(r chi.Router) {
				r.Use(authMiddleware.WithAuth)

				r.Get("/me", authHandler.GetMe)
				r.Get("/orders", usersHandler.List)
				r.Post("/orders", usersHandler.List)
				r.Get("/balance", usersHandler.List)
				r.Post("/balance/withdraw", usersHandler.List)
				r.Get("/withdrawals", usersHandler.List)
			})
		})
	})

	return appContainer.Router
}
