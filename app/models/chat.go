package models

const (
	CollectionChat = "chat"
)

type Chat struct {
	ID       int64  `json:"chat_id,omitempty" bson:"chat_id,omitempty"`
	Title    string `json:"title,omitempty" bson:"title,omitempty"`
	Username string `json:"username,omitempty" bson:"username,omitempty"`
	Type     string `json:"type,omitempty" bson:"type,omitempty"`
}
