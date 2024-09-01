package producers

import (
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/google/uuid"
)

// In Go, there is no direct equivalent to Java's `interceptor.classes`. So, we can define a custom interface
// and implement it in a struct that will be used as an interceptor.
type ProducerInterceptor interface {
	OnMessage(message *kafka.Message)
	OnAcknowledgement(message *kafka.Message)
}

type DefaultProducerInterceptor struct{}

func (dpi *DefaultProducerInterceptor) OnMessage(message *kafka.Message) {
	var traceId = uuid.New().String()
	message.Headers = append(message.Headers, kafka.Header{
		Key:   "kinactionTraceId",
		Value: []byte(traceId),
	})
	fmt.Printf("[DefaultProducerInterceptor] Added `kinactionTraceId` header to message: %s\n", traceId)
}

func (dpi *DefaultProducerInterceptor) OnAcknowledgement(message *kafka.Message) {
	if message.TopicPartition.Error != nil {
		fmt.Printf("[DefaultProducerInterceptor] Failed to deliver message: %v\n", message.TopicPartition)
	} else {
		fmt.Printf(
			"[DefaultProducerInterceptor] Produced event to topic %s: key = %-10s value = %v\n",
			*message.TopicPartition.Topic, string(message.Key), message.Value,
		)
	}
}
