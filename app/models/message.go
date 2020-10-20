package models

import (
	"fmt"
	"strings"

	"telegram-bot/pkg/utils"
)

const (
	CollectionMessage = "messages"
)

type Message struct {
	Model  `json:",inline" bson:",inline"`
	Code   string      `json:"code,omitempty" bson:"code,omitempty"`
	Action Action      `json:"action,omitempty" bson:"action,omitempty"`
	Data   interface{} `json:"data,omitempty" bson:"data,omitempty"`
}

func (m *Message) GetContent() string {
	if m.Code != "" && strings.Contains(m.Action.Content, "%s") {
		return fmt.Sprintf(m.Action.Content, m.Code)
	}
	return m.Action.Content
}

func (m *Message) GetFullContent() string {
	return m.GetContent() + "\n" + utils.Jsonify(m.Data)
}
