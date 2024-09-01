package consumers

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	customSerde "github.com/karpov-kir/kafka-in-action/serde"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

var interceptors []ConsumerInterceptor = []ConsumerInterceptor{&DefaultConsumerInterceptor{}}

func ConsumeAlertData() {
	// `value.serializer` and `key.serializer` are not supported in the Go client
	// and the respective serdes need to be invoked manually (in this example just by using
	// casting of a byte array to a string for the value and `AlertKeySerde` for the key).
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":       "localhost:9091,localhost:9092,localhost:9093",
		"enable.auto.commit":      "true",
		"group.id":                "kinaction_alertinterceptor",
		"auto.commit.interval.ms": "1000",
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create consumer: %s\n", err)
		os.Exit(1)
	}

	topic := "kinaction_alert"

	err = consumer.SubscribeTopics([]string{topic}, nil)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to subscribe to topics: %s\n", err)
		os.Exit(1)
	}

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	// A signal handler or similar could be used to set this to false to break the loop.
	run := true

	keySerde := customSerde.NewAlertKeySerde()

	for run {
		select {
		case caughtSignal := <-sigchan:
			fmt.Printf("Caught signal %v: terminating\n", caughtSignal)
			run = false
		default:
			event := consumer.Poll(100)
			if event == nil {
				continue
			}

			switch typedEvent := event.(type) {
			case *kafka.Message:
				for _, interceptor := range interceptors {
					interceptor.onConsume(typedEvent)
				}

				fmt.Printf(
					"Message (on %s, offset %d) %s: %s\n",
					typedEvent.TopicPartition,
					typedEvent.TopicPartition.Offset,
					keySerde.Deserialize(topic, typedEvent.Key),
					string(typedEvent.Value),
				)
			case kafka.Error:
				// Errors should generally be considered informational, the client will try to
				// automatically recover.
				fmt.Fprintf(os.Stderr, "%% Error: %v: %v\n", typedEvent.Code(), typedEvent)
			case kafka.OffsetsCommitted:
				for _, interceptor := range interceptors {
					interceptor.onCommit(typedEvent)
				}

				// There is no equivalent to Java's `consumer.commitAsync` in the Go client.
				// `consumer.CommitOffsets` is a blocking call and the only way to handle committed offsets
				// is to listen for the `OffsetsCommitted` event as far as I can see.
				fmt.Printf("%% Offsets committed: %v\n", typedEvent)
			default:
				fmt.Printf("Ignored %v\n", typedEvent)
			}
		}
	}

	fmt.Println("Closing consumer")
	consumer.Close()
}
