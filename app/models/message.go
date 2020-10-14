package models

type Message struct {
	Model  `json:",inline" bson:",inline"`
	Action string      `json:"name,omitempty" bson:"name,omitempty"`
	Data   interface{} `json:"data,omitempty" bson:"data,omitempty"`
}
