package worker

import (
	"context"
	"log"

	kafka_infra "github.com/toji-dev/go-shop/internal/pkg/infra/kafka-infra"
	"github.com/toji-dev/go-shop/internal/services/order-service/internal/config"
	"github.com/toji-dev/go-shop/internal/services/order-service/internal/usecase"
)

type KafkaConsumer struct {
	config       *config.Config
	orderUsecase usecase.OrderUsecase

	refundPaymentConsumer kafka_infra.Consumer
}

func NewKafkaConsumer(
	cfg *config.Config,
	orderUsecase usecase.OrderUsecase,
) *KafkaConsumer {
	consumer := &KafkaConsumer{
		config: cfg,
	}
	consumer.initKafkaConsumer()
	return consumer
}

func (sc *KafkaConsumer) initKafkaConsumer() {
	sc.refundPaymentConsumer = kafka_infra.NewConsumer(
		sc.config.Kafka.Brokers,
		"payment_events",
		"order-service-group",
	)
	log.Println("Kafka consumer initialized for topic 'payment_events'")
}

func (ks *KafkaConsumer) StartAllKafkaConsumer() {
	log.Println("Starting Kafka consumer...")
	ks.refundPaymentConsumer.Start(context.Background(), ks.orderUsecase.HandleRefundSucceededEvent)
}
