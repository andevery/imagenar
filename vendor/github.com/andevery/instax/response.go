package instax

import (
	"encoding/json"
)

type Response struct {
	Data json.RawMessage `json:"data"`
	Meta struct {
		Code         int    `json:"code"`
		ErrorMessage string `json:"error_message"`
		ErrorType    string `json:"error_type"`
	} `json:"meta"`
}
