package handlers

import (
	"errors"
	"gorm.io/gorm"
	"net/http"
	"slices"
	"ya41-56/internal/accrual/models"
	"ya41-56/internal/shared/httputil"
	models2 "ya41-56/internal/shared/models"
	"ya41-56/internal/shared/repositories"
	"ya41-56/internal/shared/response"
)

type RewardsHandler struct {
	Repo repositories.Repository[models.RewardMechanic]
}

func NewRewardsHandler(repo repositories.Repository[models.RewardMechanic]) *RewardsHandler {
	return &RewardsHandler{
		Repo: repo,
	}
}

type CreateMechanicRequest struct {
	Match      string             `json:"match"`
	Reward     float32            `json:"reward"`
	RewardType models2.RewardType `json:"reward_type"`
}

func (h *RewardsHandler) CreateRewards(w http.ResponseWriter, r *http.Request) {
	var req CreateMechanicRequest
	if err := httputil.ParseJSON(r, &req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	if req.Match == "" || (req.RewardType != "pt" && req.RewardType != "%") {
		response.Error(w, http.StatusBadRequest, "invalid mechanic data")
		return
	}

	reward := models.RewardMechanic{
		Match:      req.Match,
		Reward:     req.Reward,
		RewardType: req.RewardType,
	}

	if reward.Reward < 0 {
		response.Error(w, http.StatusBadRequest, "invalid mechanic data")
		return
	}

	if !slices.Contains(
		[]models2.RewardType{
			models2.RewardTypePercent,
			models2.RewardTypePoints,
		},
		req.RewardType,
	) {
		response.Error(w, http.StatusBadRequest, "invalid reward type")
		return
	}

	exist, err := h.Repo.FindByField(r.Context(), "match", req.Match)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			response.Error(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
	if exist.ID > 0 {
		response.Error(w, http.StatusConflict, "reward already exists")
		return
	}

	if err := h.Repo.Create(r.Context(), &reward); err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, map[string]uint{
		"id": reward.ID,
	})
}
