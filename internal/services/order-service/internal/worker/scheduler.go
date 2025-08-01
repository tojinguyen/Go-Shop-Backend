package worker

import (
	"log"

	"github.com/robfig/cron/v3"
	dependency_container "github.com/toji-dev/go-shop/internal/services/order-service/internal/dependency-container"
	"github.com/toji-dev/go-shop/internal/services/order-service/internal/usecase"
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

	// Lấy các dependencies cần thiết từ container
	orderRepo := s.container.GetOrderRepository()
	productAdapter := s.container.GetProductServiceAdapter()

	// 1. Đăng ký Job: Reconcile Pending Orders
	orderReconciler := usecase.NewOrderReconciler(orderRepo, productAdapter)

	// Thêm công việc với lịch trình. Ví dụ: chạy 5 phút một lần.
	_, err := s.cron.AddFunc("@every 5m", orderReconciler.ReconcilePendingOrders)
	if err != nil {
		// Dùng Fatalf ở đây vì nếu job không đăng ký được thì là lỗi nghiêm trọng.
		log.Fatalf("[Scheduler] FATAL: Could not register 'ReconcilePendingOrders' job: %v", err)
	}
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
