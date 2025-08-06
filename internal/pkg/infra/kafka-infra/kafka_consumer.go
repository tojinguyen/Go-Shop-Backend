package kafka_infra

import (
	"context"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type MessageHandler func(ctx context.Context, key, value []byte) error

type Consumer interface {
	Start(ctx context.Context, handler MessageHandler)
	Close() error
}

type consumer struct {
	reader  *kafka.Reader
	topic   string
	groupID string
}

// NewConsumer tạo một consumer mới.
func NewConsumer(brokers []string, topic, groupID string) Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:         brokers,
		Topic:           topic,
		GroupID:         groupID,
		MinBytes:        10e3, // 10KB
		MaxBytes:        10e6, // 10MB
		MaxWait:         10 * time.Second,
		ReadLagInterval: -1,
	})

	log.Printf("Kafka consumer created for topic '%s' with group ID '%s'", topic, groupID)
	return &consumer{
		reader:  reader,
		topic:   topic,
		groupID: groupID,
	}
}

func (c *consumer) Start(ctx context.Context, handler MessageHandler) {
	log.Printf("Starting Kafka consumer for topic '%s'...", c.topic)
	for {
		select {
		case <-ctx.Done():
			log.Printf("Context cancelled. Stopping consumer for topic '%s'.", c.topic)
			return
		default:
			// FetchMessage là một lệnh blocking, nó sẽ đợi cho đến khi có message mới hoặc context bị cancel.
			m, err := c.reader.FetchMessage(ctx)
			if err != nil {
				// Nếu context bị cancel, err sẽ là context.Canceled, vòng lặp sẽ kết thúc ở lần kiểm tra select tiếp theo.
				if err == context.Canceled {
					return
				}
				log.Printf("ERROR: Could not fetch message from topic '%s': %v", c.topic, err)
				continue // Bỏ qua lỗi và thử lại
			}

			log.Printf("Message received on topic %q, partition %d, offset %d, key: %s", m.Topic, m.Partition, m.Offset, string(m.Key))

			// Xử lý message bằng handler được cung cấp
			handleErr := handler(ctx, m.Key, m.Value)

			if handleErr != nil {
				log.Printf("ERROR: Failed to handle message for key '%s' on topic '%s'. Will not commit offset. Error: %v", string(m.Key), c.topic, handleErr)
			} else {
				// Commit message offset nếu xử lý thành công
				if err := c.reader.CommitMessages(ctx, m); err != nil {
					log.Printf("ERROR: Failed to commit message for key '%s' on topic '%s': %v", string(m.Key), c.topic, err)
				} else {
					log.Printf("Successfully committed message for key '%s' on topic '%s'", string(m.Key), c.topic)
				}
			}
		}
	}
}

// Close đóng Kafka reader.
func (c *consumer) Close() error {
	log.Printf("Closing Kafka consumer for topic '%s'...", c.topic)
	return c.reader.Close()
}
