package serde

import (
	"github.com/karpov-kir/kafka-in-action/models"
)

type AlertKeySerde struct{}

func (s *AlertKeySerde) Serialize(topic string, alert *models.Alert) []byte {
	return []byte(alert.StageID)
}

// Returns `Alert.StageID`.
func (s *AlertKeySerde) Deserialize(topic string, value []byte) string {
	return string(value)
}

func NewAlertKeySerde() *AlertKeySerde {
	return &AlertKeySerde{}
}
