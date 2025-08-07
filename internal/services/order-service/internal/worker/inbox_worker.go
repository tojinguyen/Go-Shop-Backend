package worker

import (
	"context"
	"log"
	"time"

	"github.com/toji-dev/go-shop/internal/services/order-service/internal/usecase"
)

type InboxWorker struct {
	inboxEventUsecase usecase.InboxEventUseCase
	ctx               context.Context
	cancel            context.CancelFunc
}

func NewInboxWorker(inboxEventUsecase usecase.InboxEventUseCase) *InboxWorker {
	ctx, cancel := context.WithCancel(context.Background())
	return &InboxWorker{
		inboxEventUsecase: inboxEventUsecase,
		ctx:               ctx,
		cancel:            cancel,
	}
}

func (w *InboxWorker) Start() {
	log.Println("[InboxWorker] Starting inbox processing workers...")

	// Start processing pending events
	go w.processPendingEventsWorker()

	// Start retrying failed events
	go w.retryFailedEventsWorker()

	// Start cleanup worker
	go w.cleanupWorker()

	// Start stats logging
	go w.statsWorker()

	log.Println("[InboxWorker] All inbox workers started successfully")
}

func (w *InboxWorker) Stop() {
	log.Println("[InboxWorker] Stopping inbox workers...")
	w.cancel()
}

// processPendingEventsWorker - Process pending events every 5 seconds
func (w *InboxWorker) processPendingEventsWorker() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	log.Println("[InboxWorker] Started pending events processor")

	for {
		select {
		case <-w.ctx.Done():
			log.Println("[InboxWorker] Pending events processor stopped")
			return
		case <-ticker.C:
			err := w.inboxEventUsecase.ProcessPendingInboxEvents(w.ctx)
			if err != nil {
				log.Printf("[InboxWorker] Error processing pending events: %v", err)
			}
		}
	}
}

// retryFailedEventsWorker - Retry failed events every 30 seconds
func (w *InboxWorker) retryFailedEventsWorker() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	log.Println("[InboxWorker] Started failed events retry processor")

	for {
		select {
		case <-w.ctx.Done():
			log.Println("[InboxWorker] Failed events retry processor stopped")
			return
		case <-ticker.C:
			err := w.inboxEventUsecase.RetryFailedInboxEvents(w.ctx)
			if err != nil {
				log.Printf("[InboxWorker] Error retrying failed events: %v", err)
			}
		}
	}
}

// cleanupWorker - Cleanup old events daily
func (w *InboxWorker) cleanupWorker() {
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	log.Println("[InboxWorker] Started cleanup processor")

	// Run cleanup immediately on start
	if err := w.inboxEventUsecase.CleanupOldEvents(w.ctx); err != nil {
		log.Printf("[InboxWorker] Error in initial cleanup: %v", err)
	}

	for {
		select {
		case <-w.ctx.Done():
			log.Println("[InboxWorker] Cleanup processor stopped")
			return
		case <-ticker.C:
			err := w.inboxEventUsecase.CleanupOldEvents(w.ctx)
			if err != nil {
				log.Printf("[InboxWorker] Error cleaning up old events: %v", err)
			}
		}
	}
}

// statsWorker - Log inbox statistics every 5 minutes
func (w *InboxWorker) statsWorker() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	log.Println("[InboxWorker] Started stats logger")

	for {
		select {
		case <-w.ctx.Done():
			log.Println("[InboxWorker] Stats logger stopped")
			return
		case <-ticker.C:
			stats, err := w.inboxEventUsecase.GetInboxStats(w.ctx)
			if err != nil {
				log.Printf("[InboxWorker] Error getting inbox stats: %v", err)
			} else {
				log.Printf("[InboxStats] Pending: %d, Processed: %d, Failed: %d, Total: %d",
					stats.PendingCount, stats.ProcessedCount, stats.FailedCount, stats.TotalCount)
			}
		}
	}
}
