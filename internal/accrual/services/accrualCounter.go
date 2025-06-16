package services

import (
	"math"
	"strings"
	"ya41-56/internal/accrual/models"
)

type AccrualService interface {
	CalculateAccrual(goods []models.Good, mechanics []models.RewardMechanic) float32
}

type AccrualCounter struct {
}

func NewAccrualCounter() AccrualService {
	return &AccrualCounter{}
}

func (h *AccrualCounter) CalculateAccrual(goods []models.Good, mechanics []models.RewardMechanic) float32 {
	var total float32

	for _, g := range goods {
		for _, m := range mechanics {
			if strings.Contains(strings.ToLower(g.Description), strings.ToLower(m.Match)) {
				switch m.RewardType {
				case "%":
					total += float32(g.Price) * m.Reward / 100
				case "pt":
					total += m.Reward
				}
			}
		}
	}

	return float32(math.Round(float64(total)*100) / 100)
}
