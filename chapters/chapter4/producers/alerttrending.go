package producers

import (
	"fmt"
	"os"

	"github.com/karpov-kir/kafka-in-action/models"
	customSerde "github.com/karpov-kir/kafka-in-action/serde"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func ProduceAlertTrendingData() {
	// `value.serializer` and `key.serializer` are not supported in the Go client
	// and the respective serdes need to be invoked manually (in this example just by using
	// casting of a string to a byte array for the value and `AlertKeySerde` for the key).
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9091,localhost:9092,localhost:9093",
	})

	if err != nil {
		fmt.Printf("Failed to create producer: %s\n", err)
		os.Exit(1)
	}

	topic := "kinaction_alerttrend"

	alert := models.Alert{
		AlertID:      0,
		StageID:      "Stage 0",
		AlertLevel:   "CRITICAL",
		AlertMessage: "Stage 0 stopped",
	}

	keySerde := customSerde.NewAlertKeySerde()

	// Go-routine to handle message delivery reports and
	// possibly other event types (errors, stats, etc).
	// It should be the same as the `AlertCallback` in the book.
	go func() {
		for event := range producer.Events() {
			switch typedEvent := event.(type) {
			case *kafka.Message:
				if typedEvent.TopicPartition.Error != nil {
					fmt.Printf("Failed to deliver message: %v\n", typedEvent.TopicPartition)
				} else {
					fmt.Printf(
						"Produced event to topic %s: key = %-10s value = %v\n",
						*typedEvent.TopicPartition.Topic, string(typedEvent.Key), alert,
					)
				}
			default:
				fmt.Printf("Ignored event: %s\n", event)
			}
		}
	}()

	err = producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            keySerde.Serialize(topic, &alert),
		Value:          []byte(alert.AlertMessage),
	}, nil)

	if err != nil {
		fmt.Printf("Produce failed: %v\n", err)
		os.Exit(1)
	}

	// Wait for all messages to be delivered.
	producer.Flush(15 * 1000)
	producer.Close()
}
