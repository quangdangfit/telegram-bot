package utils

import "time"

type Logs struct {
	Message    string `json:"status,omitempty" bson:"status,omitempty"`
	HandleTime string `json:"handle_time,omitempty" bson:"handle_time,omitempty"`
}

func ParseError(err error) *Logs {
	logObject := Logs{
		Message:    err.Error(),
		HandleTime: time.Now().UTC().Format(time.RFC3339Nano),
	}
	return &logObject
}
