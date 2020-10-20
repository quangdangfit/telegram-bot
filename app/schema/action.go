package schema

type Action struct {
	ID      string `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	ChatID  int64  `json:"chat_id,omitempty"`
	Content string `json:"content,omitempty"`
	Status  bool   `json:"status,omitempty"`

	CreatedTime string `json:"created_time"`
	UpdatedTime string `json:"updated_time"`
}

type ActionQueryParam struct {
	Name string `json:"name,omitempty"`
}

type ActionCreateParam struct {
	Name       string  `json:"name,omitempty" validate:"required"`
	ChatID     []int64 `json:"chat_id,omitempty" validate:"required"`
	Content    string  `json:"content,omitempty" validate:"required"`
	UrlConfirm string  `json:"url_confirm,omitempty"`
	UrlReject  string  `json:"url_reject,omitempty"`
}
