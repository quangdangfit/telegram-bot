package schema

type Action struct {
	ID     string `json:"id,omitempty"`
	Name   string `json:"name,omitempty"`
	ChatID int64  `json:"chat_id,omitempty"`
	Status bool   `json:"status,omitempty"`

	CreatedTime string `json:"created_time"`
	UpdatedTime string `json:"updated_time"`
}

type ActionQueryParam struct {
	Name string `json:"name,omitempty"`
}
