package models

const (
	CollectionInMessage = "in_messages"

	InMessageStatusReceived    = "received"
	InMessageStatusSuccess     = "success"
	InMessageStatusWaitRetry   = "wait_retry"
	InMessageStatusWorking     = "working"
	InMessageStatusFailed      = "failed"
	InMessageStatusInvalid     = "invalid"
	InMessageStatusWaitPrevMsg = "wait_prev_msg"
	InMessageStatusCanceled    = "canceled"
)

type InMessage struct {
	Model      `json:",inline" bson:",inline"`
	RoutingKey RoutingKey    `json:"routing_key,omitempty" bson:"routing_key,omitempty"`
	Payload    interface{}   `json:"payload,omitempty" bson:"payload,omitempty"`
	Status     string        `json:"status,omitempty" bson:"status,omitempty"`
	Logs       []interface{} `json:"logs,omitempty" bson:"logs,omitempty"`
	Attempts   uint          `json:"attempts" bson:"attempts"`
	Headers    `json:",inline" bson:",inline"`
}

type Headers struct {
	OriginCode  string `json:"origin_code,omitempty" bson:"origin_code,omitempty"`
	OriginModel string `json:"origin_model,omitempty" bson:"origin_model,omitempty"`
	Query       string `json:"query,omitempty" bson:"query,omitempty"`
	ExternalID  string `json:"external_id,omitempty" bson:"external_id,omitempty"`
	PublishTime string `json:"publish_time,omitempty" bson:"publish_time,omitempty"`
}
