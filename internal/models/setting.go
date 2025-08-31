package models

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type SettingKey string

const (
	SettingKeyPrompts          SettingKey = "prompts"
	SettingKeyBotName          SettingKey = "bot_name"
	SettingConversationTimeout SettingKey = "conversation_timeout"
	SettingKeyIntroduceMessage SettingKey = "introduce_message"

	PrefixCacheKey string = "setting_chatbot"
)

func MakeSettingKey(key string) string {
	return fmt.Sprintf("%s:%s", PrefixCacheKey, key)
}

var SettingMap = map[SettingKey]string{
	SettingKeyPrompts:          "Default prompts for the chatbot",
	SettingKeyBotName:          "Name of the chatbot",
	SettingConversationTimeout: "Timeout for the conversation in minutes",
	SettingKeyIntroduceMessage: "Introduce message",
}

type SettingModel struct {
	Id        bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Key       string        `bson:"key" json:"key"`
	Value     string        `bson:"value" json:"value"`
	CreatedAt time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time     `bson:"updated_at" json:"updated_at"`
	DeletedAt time.Time     `bson:"deleted_at,omitempty" json:"deleted_at,omitempty" go.mongodb.org/mongo-driver/bson.D`
}
