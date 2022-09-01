package rpc

type Request struct {
	Id      string `json:"id"`
	Method  string `json:"method"`
	Params  any    `json:"params,omitempty"`
	JsonRpc string `json:"jsonrpc"`
}
