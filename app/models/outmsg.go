package models

import "net/http"

const (
	CollectionOutMessage = "out_messages"

	OutMessageStatusWait     = "wait"
	OutMessageStatusSent     = "sent"
	OutMessageStatusSentWait = "sent_wait"
	OutMessageStatusFailed   = "failed"
	OutMessageStatusCanceled = "canceled"
	OutMessageStatusInvalid  = "invalid"
)

type OutMessage struct {
	Model       `json:",inline" bson:",inline"`
	RoutingKey  string            `json:"routing_key,omitempty" bson:"routing_key,omitempty"`
	Payload     interface{}       `json:"payload,omitempty" bson:"payload,omitempty"`
	OriginCode  string            `json:"origin_code,omitempty" bson:"origin_code,omitempty"`
	OriginModel string            `json:"origin_model,omitempty" bson:"origin_model,omitempty"`
	Status      string            `json:"status,omitempty" bson:"status,omitempty"`
	Logs        []interface{}     `json:"logs,omitempty" bson:"logs,omitempty"`
	Headers     http.Header       `json:"headers,omitempty" bson:"headers,omitempty"`
	Query       map[string]string `json:"query,omitempty" bson:"query,omitempty"`
	ExternalID  string            `json:"external_id,omitempty" bson:"external_id,omitempty"`
}
