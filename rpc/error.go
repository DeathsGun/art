package rpc

import "fmt"

type Error struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s. Code: %d", e.Message, e.Code)
}
