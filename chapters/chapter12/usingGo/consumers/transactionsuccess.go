package consumers

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde/avro"
	"github.com/karpov-kir/kafka-in-action/models"
	"github.com/karpov-kir/kafka-in-action/processors"
)

func ConsumeTransactionSuccess() {
	// `value.serializer` and `key.serializer` are not supported in the Go client
	// and the required available confluent serdes need to instantiated
	// manually; they get registered in the schema registry under the hood
	// via the `client` inside `schemaregistry.serde.GetID`.
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":       "localhost:9091,localhost:9092,localhost:9093",
		"group.id":                "transaction-success",
		"enable.auto.commit":      "true",
		"auto.commit.interval.ms": "1000",
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create consumer: %s\n", err)
		os.Exit(1)
	}

	client, err := schemaregistry.NewClient(schemaregistry.NewConfig("http://localhost:8081"))

	if err != nil {
		fmt.Printf("Failed to create schema registry client: %s\n", err)
		os.Exit(1)
	}

	deserializer, err := avro.NewSpecificDeserializer(client, serde.ValueSerde, avro.NewDeserializerConfig())

	if err != nil {
		fmt.Printf("Failed to create deserializer: %s\n", err)
		os.Exit(1)
	}

	err = consumer.SubscribeTopics([]string{processors.TransactionSuccessTopic}, nil)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to subscribe to topics: %s\n", err)
		os.Exit(1)
	}

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	// A signal handler or similar could be used to set this to false to break the loop.
	run := true

	for run {
		select {
		case caughtSignal := <-sigchan:
			fmt.Printf("Caught signal %v: terminating\n", caughtSignal)
			run = false
		default:
			message, err := consumer.ReadMessage(time.Second)
			if err == nil {
				transactionResult := models.TransactionResult{}
				err := deserializer.DeserializeInto(*message.TopicPartition.Topic, message.Value, &transactionResult)
				if err != nil {
					fmt.Printf("Failed to deserialize payload: %s\n", err)
				} else {
					fmt.Printf("%% Message on %s:\n%+v\n", message.TopicPartition, transactionResult)
				}
			} else if !err.(kafka.Error).IsTimeout() {
				// The client will automatically try to recover from all errors.
				// Timeout is not considered an error because it is raised by
				// ReadMessage in absence of messages.
				fmt.Printf("Consumer error: %v (%v)\n", err, message)
			}
		}
	}

	fmt.Println("Closing consumer")
	consumer.Close()
}
