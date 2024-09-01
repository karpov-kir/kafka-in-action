package producers

import (
	"fmt"
	"os"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde/avro"
	"github.com/karpov-kir/kafka-in-action/models"
)

func Produce() {
	// `value.serializer` and `key.serializer` are not supported in the Go client
	// and the required available confluent serdes need to instantiated
	// manually; they get registered in the schema registry under the hood
	// via the `client` inside `schemaregistry.serde.GetID`; the key is serialized just by using
	// casting of a number to a string and to a byte array.
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

	serializer, err := avro.NewSpecificSerializer(client, serde.ValueSerde, avro.NewSerializerConfig())

	if err != nil {
		fmt.Printf("Failed to create serializer: %s\n", err)
		os.Exit(1)
	}

	topic := "kinaction_schematest"

	alert := models.Alert{
		Sensor_id: 1,
		Time:      1000,
		Status:    models.AlertStatusWarning,
	}

	alertPayload, err := serializer.Serialize(topic, &alert)

	if err != nil {
		fmt.Printf("Failed to serialize payload: %s\n", err)
		os.Exit(1)
	}

	err = producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            []byte(string(alert.Sensor_id)),
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
