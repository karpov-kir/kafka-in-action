package consumers

import (
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

// In Go, there is no direct equivalent to Java's `interceptor.classes`. So, we can define a custom interface
// and implement it in a struct that will be used as an interceptor.
type ConsumerInterceptor interface {
	onConsume(message *kafka.Message)
	onCommit(event kafka.OffsetsCommitted)
}

type DefaultConsumerInterceptor struct{}

func (dci *DefaultConsumerInterceptor) onConsume(message *kafka.Message) {
	fmt.Printf(
		"[DefaultConsumerInterceptor] Consumed message with trace ID %s from topic %s: key = %-10s value = %s\n",
		getTraceIdFromHeaders(message.Headers), *message.TopicPartition.Topic, string(message.Key), string(message.Value),
	)
}

func (dci *DefaultConsumerInterceptor) onCommit(event kafka.OffsetsCommitted) {
	fmt.Printf("[DefaultConsumerInterceptor] Committed offsets: %v\n", event)
}

func getTraceIdFromHeaders(headers []kafka.Header) string {
	for _, header := range headers {
		if string(header.Key) == "kinactionTraceId" {
			return string(header.Value)
		}
	}

	return ""
}
