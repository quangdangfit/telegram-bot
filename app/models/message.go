package models

const (
	CollectionMessage = "messages"
)

type Message struct {
	Model  `json:",inline" bson:",inline"`
	Code   string      `json:"code,omitempty" bson:"code,omitempty"`
	Action string      `json:"name,omitempty" bson:"name,omitempty"`
	Data   interface{} `json:"data,omitempty" bson:"data,omitempty"`
}
