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
	"github.com/stretchr/testify/require"
)

func TestLoginSuccess(t *testing.T) {
	db := setupTestDB(t)
	repo := repositories.NewGormRepository[models.User](db)
	ctx := context.Background()

	// Create test user
	passwordHash, err := passwordHash("password")
	require.NoError(t, err)
	testUser := models.User{
		Login:        "testuser",
		PasswordHash: passwordHash,
		Status:       models.UserStatusActive,
	}
	err = repo.Create(ctx, &testUser)
	require.NoError(t, err)
	require.NotZero(t, testUser.ID)

	// Login
	authService := services.NewAuthService(repo, services.NewTokenService("secretKey", 1*time.Hour))
	token, err := authService.Login(ctx, testUser.Login, "password")
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestRegisterSuccess(t *testing.T) {
	db := setupTestDB(t)
	repo := repositories.NewGormRepository[models.User](db)
	ctx := context.Background()

	newUser := models.User{
		Login:    "testuser",
		Password: "password",
	}
	existUser, err := repo.FindByField(ctx, "login", newUser.Login)
	if err == nil {
		repo.Delete(ctx, strconv.Itoa(int(existUser.ID)))
	}

	authService := services.NewAuthService(repo, services.NewTokenService("secretKey", 1*time.Hour))
	token, err := authService.Register(ctx, &newUser)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}
