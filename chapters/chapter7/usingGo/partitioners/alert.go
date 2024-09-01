package partitioners

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/karpov-kir/kafka-in-action/models"
)

const PartitionForCriticalAlert = 0

func FindPartitionForAlert(producer *kafka.Producer, topic string, alert *models.Alert) (int32, error) {
	if isCriticalLevel(alert) {
		// If the alert is critical, we want to send it to the first partition.
		// Some other logic could be implemented here.
		return PartitionForCriticalAlert, nil
	}

	metadata, err := producer.GetMetadata(&topic, false, 5000)
	if err != nil {
		return -1, fmt.Errorf("failed to get metadata: %v", err)
	}

	partitionCount := len(metadata.Topics[topic].Partitions)
	randomPartition := rand.Intn(partitionCount)

	return int32(randomPartition), nil
}

func isCriticalLevel(alert *models.Alert) bool {
	return strings.Contains(strings.ToUpper(alert.AlertLevel), "CRITICAL")
}
