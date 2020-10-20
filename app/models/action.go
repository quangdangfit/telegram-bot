package models

const (
	CollectionAction = "actions"
)

type Action struct {
	Model      `json:",inline" bson:",inline"`
	Name       string  `json:"name,omitempty" bson:"name,omitempty"`
	ChatID     []int64 `json:"chat_id,omitempty" bson:"chat_id,omitempty"`
	Content    string  `json:"content,omitempty" bson:"content,omitempty"`
	UrlConfirm string  `json:"url_confirm,omitempty" bson:"url_confirm,omitempty"`
	UrlReject  string  `json:"url_reject,omitempty" bson:"url_reject,omitempty"`
	Status     bool    `json:"status,omitempty" bson:"status,omitempty"`
}
