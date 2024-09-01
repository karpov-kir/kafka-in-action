package producers

import (
	"fmt"
	"os"

	"github.com/karpov-kir/kafka-in-action/models"
	"github.com/karpov-kir/kafka-in-action/partitioners"
	customSerde "github.com/karpov-kir/kafka-in-action/serde"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde/jsonschema"
)

func ProduceAlertData() {
	// `value.serializer` and `key.serializer` are not supported in the Go client
	// and the respective serdes need to be invoked manually.
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9091,localhost:9092,localhost:9093",
	})

	if err != nil {
		fmt.Printf("Failed to create producer: %s\n", err)
		os.Exit(1)
	}

	client, err := schemaregistry.NewClient(schemaregistry.NewConfig("http://localhost:8081"))

	if err != nil {
		fmt.Printf("Failed to create schema registry client: %s\n", err)
		os.Exit(1)
	}

	// In Go, there is no direct equivalent to Java's `Serializable` interface,
	// In Java, `Serializable` is a marker interface in the `java.io` package.
	// It is used to indicate that a class can be serialized,
	// which means converting an object into a byte stream that can be reverted back into a copy of the object.
	// So, the simplest and closest would be to use the `jsonschema.Serializer` to serialize the object.
	// Another alternative would be to manually marshal to a byte array using the `encoding/json` package
	// and then manually unmarshal it back to the object on the consumer side.
	valueSerializer, err := jsonschema.NewSerializer(client, serde.ValueSerde, jsonschema.NewSerializerConfig())

	if err != nil {
		fmt.Printf("Failed to create serializer: %s\n", err)
		os.Exit(1)
	}

	topic := "kinaction_alert"

	alert := models.Alert{
		AlertID:      1,
		StageID:      "Stage 1",
		AlertLevel:   "CRITICAL",
		AlertMessage: "Stage 1 stopped",
	}

	alertPayload, err := valueSerializer.Serialize(topic, &alert)

	if err != nil {
		fmt.Printf("Failed to serialize alert payload: %s\n", err)
		os.Exit(1)
	}

	keySerde := customSerde.NewAlertKeySerde()

	partition, err := partitioners.FindPartitionForAlert(producer, topic, &alert)

	if err != nil {
		fmt.Printf("Failed to find partition: %s\n", err)
		os.Exit(1)
	}

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
						*typedEvent.TopicPartition.Topic, string(typedEvent.Key), alertPayload,
					)
				}
			default:
				fmt.Printf("Ignored event: %s\n", event)
			}
		}
	}()

	err = producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: partition},
		Key:            keySerde.Serialize(topic, &alert),
		Value:          alertPayload,
	}, nil)

	if err != nil {
		fmt.Printf("Produce failed: %v\n", err)
		os.Exit(1)
	}

	// Wait for all messages to be delivered.
	producer.Flush(15 * 1000)
	producer.Close()
}
