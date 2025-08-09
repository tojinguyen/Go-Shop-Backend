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
	c := cron.New(cron.WithSeconds())

	return &Scheduler{
		cron:      c,
		container: dc,
	}
}

// RegisterJobs đăng ký tất cả các công việc định kỳ.
func (s *Scheduler) RegisterJobs() {
	log.Println("[Scheduler] Registering cron jobs...")
	log.Println("[Scheduler] 'ReconcilePendingOrders' job registered to run every 5 minutes.")

	paymentEventUseCase := s.container.GetPaymentEventUseCase()

	// Job 1: Xử lý các payment success đang pending để cập nhật order status
	_, err := s.cron.AddFunc("@every 10m", paymentEventUseCase.HandleSuccessPaymentPending)

	if err != nil {
		log.Fatalf("[Scheduler] FATAL: Could not register 'ReconcilePendingOrders' job: %v", err)
	}
	log.Println("[Scheduler] 'ReconcilePendingOrders' job registered to run every 5 minutes.")

	// Job 2: Xử lý các payment failed đang pending để cập nhật order status
	_, err = s.cron.AddFunc("@every 5m", paymentEventUseCase.HandleFailedPaymentPending)
	if err != nil {
		log.Fatalf("[Scheduler] FATAL: Could not register 'HandleFailedPaymentPending' job: %v", err)
	}

	// Job 3: Xử lý các yêu cầu refund đang pending
	_, err = s.cron.AddFunc("@every 1m", paymentEventUseCase.HandleRefundPaymentPending)
	if err != nil {
		log.Fatalf("[Scheduler] FATAL: Could not register 'HandleRefundPaymentPending' job: %v", err)
	}
	log.Println("[Scheduler] 'HandleRefundPaymentPending' job registered to run every 1 minute.")

	// Job 3: Publish các event refund thành công
	_, err = s.cron.AddFunc("@every 1m", paymentEventUseCase.PublishRefundSucceededEvents)
	if err != nil {
		log.Fatalf("[Scheduler] FATAL: Could not register 'PublishRefundSucceededEvents' job: %v", err)
	}
	log.Println("[Scheduler] 'PublishRefundSucceededEvents' job registered to run every 1 minute.")
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
