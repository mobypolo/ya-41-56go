package worker

import (
	"ya41-56/internal/gophermart/models"
	"ya41-56/internal/shared/repositories"
)

type FetchPool struct {
	repo    repositories.Repository[models.Order]
	fetcher *AccrualFetcher
}

func NewFetchPool(repo repositories.Repository[models.Order], accrualAddr string) *FetchPool {
	return &FetchPool{
		repo:    repo,
		fetcher: NewFetcher(repo, accrualAddr),
	}
}

func (p *FetchPool) Start() {
	p.fetcher.RunWorkers(5) // TODO вынести бы в конфиг
}

func (p *FetchPool) Stop() {
	close(p.fetcher.jobs)
}

func (p *FetchPool) Add(order *models.Order) {
	p.fetcher.Add(order)
}
