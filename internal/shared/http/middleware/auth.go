package middleware

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"ya41-56/internal/gophermart/models"
	"ya41-56/internal/shared/contextutil"
	"ya41-56/internal/shared/response"
)

type TokenParser interface {
	ParseAndValidate(context.Context, string) (string, error)
}

type UserFinder interface {
	FindByID(context.Context, string) (*models.User, error)
}

type AuthMiddleware struct {
	TokenService TokenParser
	UserService  UserFinder
}

func New(tokenParser TokenParser, userFinder UserFinder) *AuthMiddleware {
	return &AuthMiddleware{
		TokenService: tokenParser,
		UserService:  userFinder,
	}
}

func (m *AuthMiddleware) WithAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			response.Error(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
			return
		}

		userID, err := m.TokenService.ParseAndValidate(r.Context(), strings.TrimPrefix(authHeader, "Bearer "))
		if err != nil {
			response.Error(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
			return
		}

		currentUser, err := m.UserService.FindByID(r.Context(), userID)
		if err != nil {
			response.Error(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
			return
		}

		ctx := r.Context()
		ctx = contextutil.WithUserID(ctx, strconv.Itoa(int(currentUser.ID)))

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
