package response

import "time"

type Response struct {
	Status    int         `json:"status"`
	Message   string      `json:"message"`
	Errors    interface{} `json:"errors,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	TimeStamp time.Time   `json:"time_stamp"`
}
