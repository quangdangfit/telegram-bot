package models

type Model struct {
	ID          string `json:"id,omitempty" bson:"id,omitempty"`
	CreatedTime string `json:"created_time" bson:"created_time"`
	UpdatedTime string `json:"updated_time" bson:"updated_time"`
}
