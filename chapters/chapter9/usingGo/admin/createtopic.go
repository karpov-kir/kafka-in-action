package admin

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func CreateTopic() {
	bootstrapServers := "localhost:9092,localhost:9093"
	topic := "kinaction_selfserviceTopic"

	// Create adminClient new AdminClient.
	// AdminClient can also be instantiated using an existing
	// Producer or Consumer instance, see NewAdminClientFromProducer and
	// NewAdminClientFromConsumer.
	adminClient, err := kafka.NewAdminClient(&kafka.ConfigMap{"bootstrap.servers": bootstrapServers})

	if err != nil {
		fmt.Printf("Failed to create Admin client: %s\n", err)
		os.Exit(1)
	}

	// Contexts are used to abort or limit the amount of time
	// the Admin call blocks waiting for a result.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create topics on cluster.
	// Set Admin options to wait for the operation to finish (or at most 60s)
	maxDuration, err := time.ParseDuration("60s")

	if err != nil {
		os.Exit(1)
	}

	describeTopicsResult, err := adminClient.DescribeTopics(
		ctx,
		kafka.NewTopicCollectionOfTopicNames([]string{topic}),
	)

	if err != nil {
		fmt.Printf("Failed to describe topic: %v\n", err)
		os.Exit(1)
	}

	topicDescription := describeTopicsResult.TopicDescriptions[0]

	if topicDescription.Error.Code() == kafka.ErrUnknownTopicOrPart {
		fmt.Printf("Topic %s does not exist, creating...\n", topic)
	} else if topicDescription.Error.Code() != kafka.ErrNoError {
		fmt.Printf("Failed to describe topic: %v\n", topicDescription.Error)
		os.Exit(1)
	}

	createTopicsResult, err := adminClient.CreateTopics(
		ctx,
		// Multiple topics can be created simultaneously
		// by providing more TopicSpecification structs here.
		[]kafka.TopicSpecification{{
			Topic:             topic,
			NumPartitions:     2,
			ReplicationFactor: 2,
		}},
		kafka.SetAdminOperationTimeout(maxDuration))

	if err != nil {
		fmt.Printf("Failed to create topic: %v\n", err)
		os.Exit(1)
	}

	for _, topicResult := range createTopicsResult {
		fmt.Printf("Topic created: %s\n", topicResult)
	}

	adminClient.Close()
}
