package producers

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde/avro"
	"github.com/google/uuid"
	"github.com/karpov-kir/kafka-in-action/models"
	"github.com/karpov-kir/kafka-in-action/processors"
	"github.com/karpov-kir/kafka-in-action/utils"
)

func ProduceTransactionRequests() {
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
			transaction := models.Transaction{
				Guid:     uuid.New().String(),
				Account:  getRandomAccount(),
				Amount:   getRandomAmount(),
				Type:     getRandomTransactionType(),
				Currency: "CNY",
				Country:  "CN",
			}

			transactionPayload, err := serializer.Serialize(processors.TransactionRequestTopic, &transaction)

			if err != nil {
				fmt.Printf("Failed to serialize payload: %s\n", err)
				os.Exit(1)
			}

			topic := processors.TransactionRequestTopic
			err = producer.Produce(&kafka.Message{
				TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
				Key:            []byte(transaction.Account),
				Value:          transactionPayload,
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
						"Produced event to topic %s: key = %-10s value = %v\n",
						*typedEvent.TopicPartition.Topic, string(typedEvent.Key), transaction,
					)
				}
			default:
				fmt.Printf("Unknown event: %s\n", event)
				os.Exit(1)
			}
		}
	}

	producer.Close()
}

func getRandomAccount() string {
	accounts := []string{"123", "456", "789"}
	return accounts[rand.Intn(len(accounts))]
}

func getRandomTransactionType() models.TransactionType {
	types := []models.TransactionType{models.TransactionTypeDEPOSIT, models.TransactionTypeWITHDRAW}
	return types[rand.Intn(len(types))]
}

func getRandomAmount() []byte {
	return utils.DecimalToAvroBytes(rand.Float64() * 1000)
}
