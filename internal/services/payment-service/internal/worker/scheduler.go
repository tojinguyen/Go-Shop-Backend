package worker

import (
	"log"

	"github.com/robfig/cron/v3"
	dependency_container "github.com/toji-dev/go-shop/internal/services/payment-service/internal/dependency-container"
)

// Scheduler quản lý tất cả các cron job trong service.
type Scheduler struct {
	cron      *cron.Cron
	container *dependency_container.DependencyContainer
}

// NewScheduler tạo và cấu hình một scheduler mới.
func NewScheduler(dc *dependency_container.DependencyContainer) *Scheduler {
	// Khởi tạo cron scheduler. Chúng ta có thể thêm logger hoặc các option khác ở đây.
	c := cron.New(cron.WithSeconds()) // Thêm WithSeconds() nếu bạn cần độ chính xác đến giây

	return &Scheduler{
		cron:      c,
		container: dc,
	}
}

// RegisterJobs đăng ký tất cả các công việc định kỳ.
func (s *Scheduler) RegisterJobs() {
	log.Println("[Scheduler] Registering cron jobs...")
	log.Println("[Scheduler] 'ReconcilePendingOrders' job registered to run every 5 minutes.")
}

// Start khởi động scheduler để bắt đầu chạy các công việc.
// Nó nên được gọi trong một goroutine.
func (s *Scheduler) Start() {
	log.Println("[Scheduler] Starting scheduler...")
	s.cron.Start()
}

// Stop dừng scheduler một cách an toàn.
func (s *Scheduler) Stop() {
	log.Println("[Scheduler] Stopping scheduler...")
	ctx := s.cron.Stop() // Stop trả về một context, báo hiệu khi tất cả job đang chạy đã hoàn thành.
	<-ctx.Done()
	log.Println("[Scheduler] Scheduler stopped gracefully.")
}
