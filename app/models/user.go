package models

const (
	CollectionUser = "users"
)

type User struct {
	Model    `json:",inline" bson:",inline"`
	ChatID   int64  `json:"chat_id,omitempty" bson:"chat_id,omitempty"`
	Username string `json:"username,omitempty" bson:"username,omitempty"`
}
