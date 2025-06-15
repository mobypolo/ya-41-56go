package repositories_test

import (
	"context"
	"strconv"
	"testing"
	"ya41-56/internal/accrual/models"
	"ya41-56/internal/shared/repositories"

	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&models.Order{}, &models.Good{})
	require.NoError(t, err)

	return db
}

func TestOrderRepository_CRUD(t *testing.T) {
	db := setupTestDB(t)
	repo := repositories.NewGormRepository[models.Order](db)

	ctx := context.Background()
	order := &models.Order{
		Number:  "12345678903",
		Status:  "NEW",
		Accrual: 12.5,
		Goods: []models.Good{
			{Description: "Чипсы Pringles", Price: 100},
			{Description: "Coca-Cola", Price: 55},
		},
	}

	err := repo.Create(ctx, order)
	require.NoError(t, err)
	require.NotZero(t, order.ID)

	found, err := repo.FindByID(ctx, strconv.Itoa(int(order.ID)))
	require.NoError(t, err)
	require.Equal(t, order.Number, found.Number)

	order.Status = "PROCESSED"
	err = repo.Update(ctx, order)
	require.NoError(t, err)

	updated, err := repo.FindByID(ctx, strconv.Itoa(int(order.ID)))
	require.NoError(t, err)
	require.Equal(t, "PROCESSED", updated.Status)

	err = repo.Delete(ctx, strconv.Itoa(int(order.ID)))
	require.NoError(t, err)

	_, err = repo.FindByID(ctx, strconv.Itoa(int(order.ID)))
	require.Error(t, err)
	require.Equal(t, gorm.ErrRecordNotFound, err)
}
