package producers

import (
	"fmt"
	"os"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func Produce() {
	// `value.serializer` and `key.serializer` are not supported in the Go client
	// and the respective serdes need to be invoked manually (in this example just by using
	// casting of a string to a byte array).
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9091,localhost:9092,localhost:9093",
	})

	if err != nil {
		fmt.Printf("Failed to create producer: %s\n", err)
		os.Exit(1)
	}

	topic := "kinaction_helloworld"

	err = producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            nil,
		Value:          []byte("kinaction_helloworld"),
	}, nil)

	if err != nil {
		fmt.Printf("Produce failed: %v\n", err)
		os.Exit(1)
	}

	// Wait for all messages to be delivered.
	producer.Flush(15 * 1000)
	producer.Close()
}
