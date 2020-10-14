package schema

import "transport/lib/utils/paging"

type ResponsePagingResult struct {
	Paging *paging.Paging `json:"paging"`
	Data   interface{}    `json:"data"`
}

type ExternalResponse struct {
	Status     interface{} `json:"status,omitempty" bson:"status,omitempty"`
	StatusCode int         `json:"status_code,omitempty" bson:"status_code,omitempty"`
	ErrorCode  string      `json:"err_code,omitempty" bson:"err_code,omitempty"`
	Data       interface{} `json:"data,omitempty" bson:"data,omitempty"`
	Message    string      `json:"message,omitempty" bson:"message,omitempty"`
	TraceID    string      `json:"trace_id,omitempty" bson:"trace_id,omitempty"`
	Body       interface{} `json:"body,omitempty" bson:"body,omitempty"`
	Latency    int64       `json:"latency,omitempty" bson:"latency,omitempty"`
	HandleTime string      `json:"handle_time,omitempty" bson:"handle_time,omitempty"`
}
