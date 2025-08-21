package _type

import "github.com/uchupx/saceri-chatbot-api/internal/models"

type SettingUpdateRequest struct {
	Key   models.SettingKey `json:"key"`
	Value string            `json:"value"`
}
