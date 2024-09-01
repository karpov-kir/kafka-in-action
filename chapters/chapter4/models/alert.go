package models

type Alert struct {
	AlertID      int64  `json:"sensor_id"`
	StageID      string `json:"stage_id"`
	AlertLevel   string `json:"alert_level"`
	AlertMessage string `json:"alert_message"`
}
