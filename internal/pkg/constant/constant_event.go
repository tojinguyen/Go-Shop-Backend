package constant

type EventType string

const (
	EventTypeRefundSuccessed EventType = "refund_succeeded"
)

type KafkaConsumerGroupName string

const (
	KafkaConsumerGroupOrderService KafkaConsumerGroupName = "order-service-group"
)
