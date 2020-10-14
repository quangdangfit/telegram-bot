package models

type RoutingAuthType int

const (
	BasicAuthentication  RoutingAuthType = 1
	APIKeyAuthentication RoutingAuthType = 2
	CollectionRoutingKey                 = "routing_keys"
)

type RoutingKey struct {
	ID           string          `json:"id,omitempty" bson:"id,omitempty"`
	Name         string          `json:"name,omitempty" bson:"name,omitempty"`
	Group        string          `json:"group,omitempty" bson:"group,omitempty"`
	Value        uint            `json:"value,omitempty" bson:"value,omitempty"`
	APIMethod    string          `json:"api_method,omitempty" bson:"api_method,omitempty"`
	APIUrl       string          `json:"api_url,omitempty" bson:"api_url,omitempty"`
	Active       bool            `json:"active,omitempty" bson:"active,omitempty"`
	AuthType     RoutingAuthType `json:"auth_type,omitempty" bson:"auth_type,omitempty"`
	AuthKey      string          `json:"auth_key,omitempty" bson:"auth_key,omitempty"`
	ByPassStatus []interface{}   `json:"by_pass_status,omitempty" bson:"by_pass_status,omitempty"`
	RetryTimes   uint            `json:"retry_times,omitempty" bson:"retry_times,omitempty"`

	CreatedTime string `json:"created_time,omitempty" bson:"created_time,omitempty"`
	UpdatedTime string `json:"updated_time,omitempty" bson:"updated_time,omitempty"`
}
