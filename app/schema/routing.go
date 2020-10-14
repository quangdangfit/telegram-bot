package schema

import "transport/ems/app/models"

type RoutingKey struct {
	ID         string `json:"id,omitempty"`
	Name       string `json:"name,omitempty"`
	Group      string `json:"group,omitempty"`
	Value      uint   `json:"value"`
	APIMethod  string `json:"api_method,omitempty"`
	APIUrl     string `json:"api_url,omitempty"`
	Active     bool   `json:"active,omitempty"`
	RetryTimes uint   `json:"retry_times"`

	AuthType     models.RoutingAuthType `json:"auth_type"`
	AuthKey      string                 `json:"auth_key,omitempty"`
	ByPassStatus []interface{}          `json:"by_pass_status,omitempty"`

	CreatedTime string `json:"created_time,omitempty"`
	UpdatedTime string `json:"updated_time,omitempty"`
}

type RoutingKeyQueryParam struct {
	Group string `json:"group,omitempty" form:"group,omitempty"`
	Name  string `json:"name,omitempty" form:"name,omitempty"`
	Value uint   `json:"value,omitempty" form:"value,omitempty"`
	Page  int    `json:"-" form:"page,omitempty"`
	Limit int    `json:"-" form:"limit,omitempty"`
}

type RoutingKeyBodyCreateParam struct {
	Name       string                 `json:"name,omitempty" validate:"required"`
	Group      string                 `json:"group,omitempty" validate:"required"`
	Value      uint                   `json:"value,omitempty" validate:"required,gte=0"`
	APIMethod  string                 `json:"api_method,omitempty" validate:"required,oneof=GET POST PUT DELETE PATCH"`
	APIUrl     string                 `json:"api_url,omitempty" validate:"required,url"`
	AuthType   models.RoutingAuthType `json:"auth_type,omitempty" validate:"required,gte=1,lte=2"`
	AuthKey    string                 `json:"auth_key,omitempty" validate:"required"`
	RetryTimes uint                   `json:"retry_times,omitempty" validate:"gte=0"`
}

type RoutingKeyBodyUpdateParam struct {
	Name       string                 `json:"name,omitempty"`
	Group      string                 `json:"group,omitempty"`
	Value      uint                   `json:"value" validate:"gte=0"`
	APIMethod  string                 `json:"api_method,omitempty" validate:"oneof=GET POST PUT DELETE PATCH"`
	APIUrl     string                 `json:"api_url,omitempty" validate:"url"`
	AuthType   models.RoutingAuthType `json:"auth_type,omitempty" validate:"gte=1,lte=2"`
	AuthKey    string                 `json:"auth_key,omitempty"`
	RetryTimes uint                   `json:"retry_times" validate:"gte=0"`
}
