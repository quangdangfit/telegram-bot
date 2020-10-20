package schema

type Message struct {
	ID     string      `json:"id,omitempty"`
	Code   string      `json:"code,omitempty"`
	Action string      `json:"action,omitempty"`
	Data   interface{} `json:"data,omitempty"`

	CreatedTime string `json:"created_time"`
	UpdatedTime string `json:"updated_time"`
}

type MessageCreateParam struct {
	Code   string      `json:"code,omitempty" validate:"required"`
	Action string      `json:"action,omitempty" validate:"required"`
	Data   interface{} `json:"data,omitempty" validate:"required"`
}
