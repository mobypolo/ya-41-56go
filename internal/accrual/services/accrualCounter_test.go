package services_test

import (
	"testing"

	"ya41-56/internal/accrual/models"
	"ya41-56/internal/accrual/services"

	"github.com/stretchr/testify/require"
)

func TestCalculateAccrual(t *testing.T) {
	mechanics := []models.RewardMechanic{
		{Match: "Coca-Cola", Reward: 12.5, RewardType: "pt"},
		{Match: "Pringles", Reward: 10, RewardType: "%"},
	}

	goods := []models.Good{
		{Description: "Напиток Coca-Cola Zero 0.5л", Price: 55},
		{Description: "Чипсы Pringles бекон", Price: 147.5},
	}

	accrual := services.NewAccrualCounter().CalculateAccrual(goods, mechanics)

	require.Equal(t, float32(27.25), accrual, "сумма начислений должна быть 27.25")
}
