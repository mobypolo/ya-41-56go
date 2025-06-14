package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"net/http"
	"ya41-56/internal/gophermart/di"
	"ya41-56/internal/gophermart/handlers"
	sharedHandlers "ya41-56/internal/shared/handlers"
	"ya41-56/internal/shared/http/middleware"
)

func RegisterRoutes(appContainer *di.AppContainer) http.Handler {
	authMiddleware := middleware.New(*appContainer.Auth)

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

	// TODO: Нигде в методах нет реализации, чисто заглушки
	appContainer.Router.Route("/api", func(r chi.Router) {
		r.Route("/user", func(r chi.Router) {
			r.Post("/login", authHandler.Login)
			r.Post("/register", authHandler.Register)

			r.Group(func(r chi.Router) {
				r.Use(authMiddleware.IsAuthenticated)

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
