package kafka_infra

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type Producer interface {
	Publish(ctx context.Context, topic string, key string, value interface{}) error
	Close() error
}

type producer struct {
	writer *kafka.Writer
}

func NewProducer(brokers []string) Producer {
	writer := &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireAll, // Đảm bảo an toàn dữ liệu
		Async:        false,            // Gửi đồng bộ để biết chắc chắn đã gửi
		WriteTimeout: 10 * time.Second,
	}
	return &producer{writer: writer}
}

func (p *producer) Publish(ctx context.Context, topic string, key string, value interface{}) error {
	payload, err := json.Marshal(value)
	if err != nil {
		log.Printf("Failed to marshal kafka message value: %v", err)
		return err
	}

	msg := kafka.Message{
		Topic: topic,
		Key:   []byte(key),
		Value: payload,
	}

	err = p.writer.WriteMessages(ctx, msg)
	if err != nil {
		log.Printf("Failed to write kafka message to topic %s: %v", topic, err)
		return err
	}
	log.Printf("Successfully published message to topic %s with key %s", topic, key)
	return nil
}

func (p *producer) Close() error {
	log.Println("Closing Kafka producer...")
	return p.writer.Close()
}
