package models

const (
	CollectionMessage = "messages"
)

type Message struct {
	Model  `json:",inline" bson:",inline"`
	Code   string      `json:"code,omitempty" bson:"code,omitempty"`
	Action Action      `json:"action,omitempty" bson:"action,omitempty"`
	Data   interface{} `json:"data,omitempty" bson:"data,omitempty"`
}
