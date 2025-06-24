package worker

import (
	"context"
	"time"
	"ya41-56/internal/accrual/models"
	"ya41-56/internal/accrual/services"
	"ya41-56/internal/shared/logger"
	"ya41-56/internal/shared/repositories"

	"go.uber.org/zap"
)

type Pool struct {
	ordersRepo  repositories.Repository[models.Order]
	rewardsRepo repositories.Repository[models.RewardMechanic]
	tasks       chan OrderTask
	done        chan struct{}
}

func NewPool(orders repositories.Repository[models.Order], rewards repositories.Repository[models.RewardMechanic], workers int) *Pool {
	p := &Pool{
		ordersRepo:  orders,
		rewardsRepo: rewards,
		tasks:       make(chan OrderTask, 1000),
		done:        make(chan struct{}),
	}

	for i := 0; i < workers; i++ {
		go p.worker(i)
	}

	return p
}

func (p *Pool) Add(task OrderTask) {
	select {
	case p.tasks <- task:
	default:
		logger.L().Warn("task queue is full, dropping task", zap.Uint("order_id", task.OrderID))
	}
}

func (p *Pool) Stop() {
	close(p.tasks)
	<-p.done
}

func (p *Pool) worker(id int) {
	for task := range p.tasks {
		logger.L().Debug("worker picked task", zap.Int("worker_id", id), zap.Uint("order_id", task.OrderID))

		ctx := context.Background()

		time.Sleep(1 * time.Second)

		orderWithGoods, err := p.ordersRepo.FindByIDWithPreloads(ctx, task.OrderID, "Goods")
		if err != nil {
			logger.L().Error("failed to preload goods", zap.Uint("id", task.OrderID), zap.Error(err))
			continue
		}

		mechanics, err := p.rewardsRepo.FindAll(ctx)
		if err != nil {
			logger.L().Error("failed to load mechanics", zap.Error(err))
			continue
		}

		accrual := services.NewAccrualCounter().CalculateAccrual(orderWithGoods.Goods, mechanics)
		orderWithGoods.Accrual = accrual
		orderWithGoods.Status = "PROCESSED"

		if err := p.ordersRepo.Update(ctx, orderWithGoods); err != nil {
			logger.L().Error("failed to update order", zap.Uint("id", task.OrderID), zap.Error(err))
		}
	}

	close(p.done)
}
