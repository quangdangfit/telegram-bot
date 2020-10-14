package schema

import "net/http"

type OutMessage struct {
	ID          string            `json:"id,omitempty"`
	RoutingKey  string            `json:"routing_key,omitempty"`
	Payload     interface{}       `json:"payload,omitempty"`
	OriginCode  string            `json:"origin_code,omitempty"`
	OriginModel string            `json:"origin_model,omitempty"`
	Status      string            `json:"status,omitempty"`
	Logs        []interface{}     `json:"logs,omitempty"`
	Headers     http.Header       `json:"headers,omitempty"`
	Query       map[string]string `json:"query,omitempty"`
	ExternalID  string            `json:"external_id,omitempty"`
	PublishTime string            `json:"publish_time,omitempty"`

	CreatedTime string `json:"created_time" bson:"created_time"`
	UpdatedTime string `json:"updated_time" bson:"updated_time"`
}

type OutMessageQueryParam struct {
	RoutingKey  string `json:"routing_key,omitempty" form:"routing_key,omitempty"`
	OriginCode  string `json:"origin_code,omitempty" form:"origin_code,omitempty"`
	OriginModel string `json:"origin_model,omitempty" form:"origin_model,omitempty"`
	Status      string `json:"status,omitempty" form:"status,omitempty"`
	CreatedTime string `json:"-" form:"created_time,omitempty"`
	Page        int    `json:"-" form:"page,omitempty"`
	Limit       int    `json:"-" form:"limit,omitempty"`
	Sort        string `json:"-" form:"sort,omitempty"`
}

type OutMessageInStatusQueryParam struct {
	Status      []string `json:"-" form:"status,omitempty"`
	CreatedDays int      `json:"-" form:"created_days,omitempty"`
	UpdatedDays int      `json:"-" form:"updated_days,omitempty"`
	Page        int      `json:"-" form:"page,omitempty"`
	Limit       int      `json:"-" form:"limit,omitempty"`
	Sort        string   `json:"-" form:"sort,omitempty"`
}

type OutMessageBodyParam struct {
	RoutingKey  string                 `json:"routing_key,omitempty" validate:"required"`
	Payload     map[string]interface{} `json:"payload,omitempty" validate:"required"`
	OriginCode  string                 `json:"origin_code,omitempty"`
	OriginModel string                 `json:"origin_model,omitempty"`
	Headers     map[string][]string    `json:"headers,omitempty"`
	Query       map[string]string      `json:"query,omitempty"`
	ExternalID  string                 `json:"external_id,omitempty"`
}

type OutMessageBodyUpdateParam struct {
	Status string `json:"status,omitempty" bson:"status,omitempty" validate:"oneof=wait canceled failed sent"`

	Attempts int                    `json:"attempts,omitempty" bson:"attempts,omitempty" validate:"gte=0"`
	Payload  map[string]interface{} `json:"payload,omitempty" bson:"payload,omitempty"`
}
