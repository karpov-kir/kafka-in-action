package producers

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/karpov-kir/kafka-in-action/partitioners"
	kafkaTestContainers "github.com/testcontainers/testcontainers-go/modules/kafka"
)

const topic = "kinaction_alert"

// There is no direct equivalent to Java's `EmbeddedKafkaCluster` in Go.
// Instead, we can use the `testcontainers-go` library to run a Kafka container with a broker.
func TestKafkaProducerPartitionerInteg(t *testing.T) {
	ctx := context.Background()

	bootstrapServers, terminate := runKafkaContainer(t, ctx)

	defer func() {
		if err := terminate(); err != nil {
			t.Fatalf("Failed to terminate container: %s", err)
		}
	}()

	ProduceAlertData(bootstrapServers)
	message := getAlertMessage(t, bootstrapServers)

	if message.TopicPartition.Partition != partitioners.PartitionForCriticalAlert {
		t.Fatalf("Expected partition %d, but got %d", partitioners.PartitionForCriticalAlert, message.TopicPartition.Partition)
	}
}

func runKafkaContainer(t *testing.T, ctx context.Context) (string, func() error) {
	kafkaContainer, err := kafkaTestContainers.Run(ctx,
		"confluentinc/confluent-local:7.5.0",
		kafkaTestContainers.WithClusterID("test-cluster"),
	)

	if err != nil {
		t.Fatalf("failed to start container: %s", err)
	}

	brokers, err := kafkaContainer.Brokers(ctx)

	if err != nil {
		t.Fatalf("failed to get brokers: %s", err)
	}

	bootstrapServers := strings.Join(brokers, ",")

	createTopic(t, bootstrapServers, ctx)

	return bootstrapServers, func() error {
		return kafkaContainer.Terminate(ctx)
	}
}

func createTopic(t *testing.T, bootstrapServers string, context context.Context) {
	adminClient, err := kafka.NewAdminClient(&kafka.ConfigMap{"bootstrap.servers": bootstrapServers})

	if err != nil {
		t.Fatalf("failed to create admin client: %s", err)
	}

	defer adminClient.Close()

	topicSpecification := kafka.TopicSpecification{
		Topic:             topic,
		NumPartitions:     5,
		ReplicationFactor: 1,
	}

	createTopicsResult, err := adminClient.CreateTopics(context, []kafka.TopicSpecification{topicSpecification})

	if err != nil {
		t.Fatalf("failed to create topic: %s", err)
	}

	if createTopicsResult[0].Error.Code() != kafka.ErrNoError {
		t.Fatalf("failed to create topic: %s", createTopicsResult[0].Error)
	}

	describeTopicsResult, err := adminClient.DescribeTopics(context, kafka.NewTopicCollectionOfTopicNames([]string{topic}))

	if err != nil {
		t.Fatalf("failed to describe topic: %s", err)
	}

	topicDescription := describeTopicsResult.TopicDescriptions[0]

	if topicDescription.Error.Code() != kafka.ErrNoError {
		t.Fatalf("failed to describe topic: %s", topicDescription.Error)
	}

	// Ensure that there are more than 1 partition to ensure that the partitioner has some partitions to work with.
	if len(topicDescription.Partitions) <= 1 {
		t.Fatalf("expected more than 1 partition, but got %d", len(topicDescription.Partitions))
	}

}

func getAlertMessage(t *testing.T, bootstrapServers string) *kafka.Message {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": bootstrapServers,
		"group.id":          "test",
		// Read the messages from the beginning.
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		t.Fatalf("Failed to create consumer: %s", err)
	}

	defer func() {
		if err := consumer.Close(); err != nil {
			t.Fatalf("Failed to close consumer: %s", err)
		}
	}()

	err = consumer.SubscribeTopics([]string{topic}, nil)

	if err != nil {
		t.Fatalf("Failed to subscribe to topics: %s", err)
	}

	message, err := consumer.ReadMessage(time.Second * 5)

	if err != nil {
		t.Fatalf("Failed to read message: %s", err)
	}

	return message
}
