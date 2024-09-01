package consumers

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func ConsumeAuditData() {
	// `value.serializer` and `key.serializer` are not supported in the Go client
	// and the respective serdes need to be invoked manually (in this example just by using
	// casting of a byte array to a string).
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  "localhost:9091,localhost:9092,localhost:9093",
		"enable.auto.commit": "false",
		"group.id":           "kinaction_group_audit",
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create consumer: %s\n", err)
		os.Exit(1)
	}

	topic := "kinaction_audit"

	err = consumer.SubscribeTopics([]string{topic}, nil)

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
				fmt.Printf("Message on %s: %s\n", message.TopicPartition, string(message.Value))
				commitOffset(consumer, *message.TopicPartition.Topic, int(message.TopicPartition.Partition), int(message.TopicPartition.Offset))
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

func commitOffset(consumer *kafka.Consumer, topic string, partition int, offset int) {
	response, err := consumer.CommitOffsets([]kafka.TopicPartition{{
		Topic:     &topic,
		Partition: int32(partition),
		Offset:    kafka.Offset(offset + 1),
	}})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to commit offset: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Partition %d offset %d committed successfully", response[0].Partition, offset)
}
