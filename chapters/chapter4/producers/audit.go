package producers

import (
	"fmt"
	"os"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func ProduceAuditData() {
	// `value.serializer` and `key.serializer` are not supported in the Go client
	// and the respective serdes need to be invoked manually (in this example just by using
	// casting of a string to a byte array for the value).
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers":                     "localhost:9091,localhost:9092,localhost:9093",
		"acks":                                  "all",
		"retries":                               3,
		"max.in.flight.requests.per.connection": 1,
	})

	if err != nil {
		fmt.Printf("Failed to create producer: %s\n", err)
		os.Exit(1)
	}

	topic := "kinaction_audit"

	err = producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            nil,
		Value:          []byte("audit event"),
	}, nil)

	if err != nil {
		fmt.Printf("Produce failed: %v\n", err)
		os.Exit(1)
	}

	// Wait until the message is delivered.
	event := <-producer.Events()

	switch typedEvent := event.(type) {
	case *kafka.Message:
		if typedEvent.TopicPartition.Error != nil {
			fmt.Printf("Failed to deliver message: %v\n", typedEvent.TopicPartition)
			os.Exit(1)
		} else {
			fmt.Printf(
				"Produced event to topic %s: key = %-10s value = %s\n",
				*typedEvent.TopicPartition.Topic, string(typedEvent.Key), string(typedEvent.Value),
			)
		}
	default:
		fmt.Printf("Unknown event: %s\n", event)
		os.Exit(1)
	}

	producer.Close()
}
