package worker

import (
	"context"
	"log"

	"github.com/toji-dev/go-shop/internal/pkg/constant"
	kafka_infra "github.com/toji-dev/go-shop/internal/pkg/infra/kafka-infra"
	"github.com/toji-dev/go-shop/internal/services/order-service/internal/config"
	"github.com/toji-dev/go-shop/internal/services/order-service/internal/usecase"
)

type KafkaConsumer struct {
	config            *config.Config
	orderUsecase      usecase.OrderUsecase      // Deprecated: for backward compatibility
	inboxEventUsecase usecase.InboxEventUseCase // New inbox-based processing

	refundPaymentConsumer kafka_infra.Consumer
}

func NewKafkaConsumer(
	cfg *config.Config,
	orderUsecase usecase.OrderUsecase,
	inboxEventUsecase usecase.InboxEventUseCase,
) *KafkaConsumer {
	consumer := &KafkaConsumer{
		config:            cfg,
		orderUsecase:      orderUsecase,
		inboxEventUsecase: inboxEventUsecase,
	}
	consumer.initKafkaConsumer()
	return consumer
}

func (sc *KafkaConsumer) initKafkaConsumer() {
	sc.refundPaymentConsumer = kafka_infra.NewConsumer(
		sc.config.Kafka.Brokers,
		string(constant.EventTypeRefundSuccessed),
		string(constant.KafkaConsumerGroupOrderService),
	)
	log.Println("Kafka consumer initialized for topic 'payment_events'")
}

func (ks *KafkaConsumer) StartAllKafkaConsumer() {
	log.Println("Starting Kafka consumer with inbox pattern...")
	// Use inbox event usecase for idempotent processing
	ks.refundPaymentConsumer.Start(context.Background(), ks.inboxEventUsecase.HandleIncomingEvent)
}
