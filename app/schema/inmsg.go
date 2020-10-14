package schema

type InMessage struct {
	ID          string        `json:"id,omitempty"`
	RoutingKey  RoutingKey    `json:"routing_key,omitempty"`
	Payload     interface{}   `json:"payload,omitempty"`
	Status      string        `json:"status,omitempty"`
	Logs        []interface{} `json:"logs,omitempty"`
	Attempts    uint          `json:"attempts"`
	OriginCode  string        `json:"origin_code,omitempty"`
	OriginModel string        `json:"origin_model,omitempty"`
	Query       string        `json:"query,omitempty"`
	ExternalID  string        `json:"external_id,omitempty"`
	PublishTime string        `json:"publish_time,omitempty" bson:"publish_time,omitempty"`

	CreatedTime string `json:"created_time"`
	UpdatedTime string `json:"updated_time"`
}

type InMessageQueryParam struct {
	RoutingKey   string `json:"routing_key.name,omitempty" form:"routing_key.name,omitempty"`
	RoutingGroup string `json:"routing_key.group,omitempty" form:"routing_key.group,omitempty"`
	RoutingValue uint   `json:"routing_key.value,omitempty" form:"routing_key.value,omitempty"`
	OriginCode   string `json:"origin_code,omitempty" form:"origin_code,omitempty"`
	OriginModel  string `json:"origin_model,omitempty" form:"origin_model,omitempty"`
	Status       string `json:"status,omitempty" form:"status,omitempty"`
	CreatedTime  string `json:"-" form:"created_time,omitempty"`
	Page         int    `json:"-" form:"page,omitempty"`
	Limit        int    `json:"-" form:"limit,omitempty"`
	Sort         string `json:"-" form:"sort,omitempty"`
}

type InMessagePreviousQueryParam struct {
	OriginCode  string `json:"origin_code,omitempty" form:"origin_code,omitempty"`
	OriginModel string `json:"origin_model,omitempty" form:"origin_model,omitempty"`
	PublishTime string `json:"-" form:"publish_time,omitempty"`
	Page        int    `json:"-" form:"page,omitempty"`
	Limit       int    `json:"-" form:"limit,omitempty"`
	Sort        string `json:"-" form:"sort,omitempty"`
}

type InMessageInStatusQueryParam struct {
	Status      []string `json:"-" form:"status,omitempty"`
	CreatedDays int      `json:"-" form:"created_days,omitempty"`
	UpdatedDays int      `json:"-" form:"updated_days,omitempty"`
	Page        int      `json:"-" form:"page,omitempty"`
	Limit       int      `json:"-" form:"limit,omitempty"`
	Sort        string   `json:"-" form:"sort,omitempty"`
}

type InMessageBodyUpdateParam struct {
	Status string `json:"status,omitempty" bson:"status,omitempty" validate:"oneof=received success wait_retry working failed invalid wait_prev_msg canceled"`

	Attempts int                    `json:"attempts" bson:"attempts" validate:"gte=0"`
	Payload  map[string]interface{} `json:"payload,omitempty" bson:"payload,omitempty"`
}
