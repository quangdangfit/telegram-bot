package schema

import "transport/lib/utils/paging"

type ResponsePagingResult struct {
	Paging *paging.Paging `json:"paging"`
	Data   interface{}    `json:"data"`
}

type ExternalResponse struct {
	Status    interface{} `json:"status,omitempty"`
	ErrorCode string      `json:"err_code"`
	Message   string      `json:"message"`
	TraceID   string      `json:"trace_id"`
}
