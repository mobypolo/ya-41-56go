package services_test

import (
	"context"
	"strconv"
	"testing"
	"time"
	"ya41-56/internal/gophermart/models"
	"ya41-56/internal/gophermart/services"
	"ya41-56/internal/shared/repositories"

	"github.com/stretchr/testify/assert"
)

func TestBuildJWTString(t *testing.T) {
	db := setupTestDB(t)
	repo := repositories.NewGormRepository[models.User](db)

	authService := services.NewAuthService(repo, services.NewTokenService("secretKey", 1*time.Hour))
	token, err := authService.TokenService.BuildJWTString(1)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestParseAndValidateSuccess(t *testing.T) {
	db := setupTestDB(t)
	repo := repositories.NewGormRepository[models.User](db)
	ctx := context.Background()

	// Create test user
	testUser := models.User{
		Login:    "testuser",
		Password: "password",
	}
	existUser, err := repo.FindByField(ctx, "login", testUser.Login)
	if err == nil {
		repo.Delete(ctx, strconv.Itoa(int(existUser.ID)))
	}

	// Login
	authService := services.NewAuthService(repo, services.NewTokenService("secretKey", 1*time.Hour))
	token, err := authService.Register(ctx, &testUser)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	userID, err := authService.ParseAndValidate(ctx, token)
	assert.NoError(t, err)
	assert.NotEqual(t, 0, userID)
}
