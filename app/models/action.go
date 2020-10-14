package models

const (
	CollectionAction = "actions"
)

type Action struct {
	Model  `json:",inline" bson:",inline"`
	Name   int64 `json:"name,omitempty" bson:"name,omitempty"`
	ChatID int64 `json:"chat_id,omitempty" bson:"chat_id,omitempty"`
	Status bool  `json:"status,omitempty" bson:"status,omitempty"`
}
