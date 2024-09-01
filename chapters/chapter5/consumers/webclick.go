package consumers

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func ConsumePromoDataForMagicCalculation() {
	// `value.serializer` and `key.serializer` are not supported in the Go client
	// and the respective serdes need to be invoked manually (in this example just by using
	// casting of a byte array to a number).
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":       "localhost:9091,localhost:9092,localhost:9093",
		"enable.auto.commit":      "true",
		"auto.commit.interval.ms": "1000",
		"group.id":                "kinaction_webconsumer",
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create consumer: %s\n", err)
		os.Exit(1)
	}

	topic := "kinaction_promo"

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
				promoValue, err := strconv.ParseFloat(string(message.Value), 64)

				if err != nil {
					fmt.Printf("Error converting message value to int32: %v\n", err)
				}

				magicValue := promoValue * 1.543

				fmt.Printf("Message on %s: %f\n", message.TopicPartition, magicValue)
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
