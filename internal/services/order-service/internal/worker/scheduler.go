package worker

import (
	"log"

	"github.com/robfig/cron/v3"
	dependency_container "github.com/toji-dev/go-shop/internal/services/order-service/internal/dependency-container"
	"github.com/toji-dev/go-shop/internal/services/order-service/internal/usecase"
)

type Scheduler struct {
	cron      *cron.Cron
	container *dependency_container.DependencyContainer
}

func NewScheduler(dc *dependency_container.DependencyContainer) *Scheduler {
	c := cron.New(cron.WithSeconds())

	return &Scheduler{
		cron:      c,
		container: dc,
	}
}

func (s *Scheduler) RegisterJobs() {
	log.Println("[Scheduler] Registering cron jobs...")

	orderRepo := s.container.GetOrderRepository()
	productAdapter := s.container.GetProductServiceAdapter()

	orderReconciler := usecase.NewOrderReconciler(orderRepo, productAdapter)

	_, err := s.cron.AddFunc("@every 5m", orderReconciler.ReconcilePendingOrders)
	if err != nil {
		log.Fatalf("[Scheduler] FATAL: Could not register 'ReconcilePendingOrders' job: %v", err)
	}
	log.Println("[Scheduler] 'ReconcilePendingOrders' job registered to run every 5 minutes.")
}

func (s *Scheduler) Start() {
	log.Println("[Scheduler] Starting scheduler...")
	s.cron.Start()
}

func (s *Scheduler) Stop() {
	log.Println("[Scheduler] Stopping scheduler...")
	ctx := s.cron.Stop()
	<-ctx.Done()
	log.Println("[Scheduler] Scheduler stopped gracefully.")
}
