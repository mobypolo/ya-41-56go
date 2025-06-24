package services

import (
	"math"
	"strings"
	"ya41-56/internal/accrual/models"
	sharedModels "ya41-56/internal/shared/models"
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
	var total float64

	for _, g := range goods {
		for _, m := range mechanics {
			if strings.Contains(strings.ToLower(g.Description), strings.ToLower(m.Match)) {
				switch m.RewardType {
				case sharedModels.RewardTypePercent:
					total += float64(g.Price) * float64(m.Reward) / 100
				case sharedModels.RewardTypePoints:
					total += float64(m.Reward)
				}
			}
		}
	}

	return float32(math.Round(total*100) / 100)
}
