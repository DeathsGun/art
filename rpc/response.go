package rpc

type Response struct {
	Id     string `json:"id"`
	Result any    `json:"result"`
	Error  *Error `json:"error,omitempty"`
}
