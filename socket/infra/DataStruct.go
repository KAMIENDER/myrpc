//+build linux darwin windows

package infra

import "github.com/vmihailenco/msgpack"

type RPCRequest struct {
	MethodName string        `json:"method_name"`
	Params     []interface{} `json:"params"`
}

func NewRPCRequest(methodName string, params []interface{}) *RPCRequest {
	return &RPCRequest{
		MethodName: methodName,
		Params:     params,
	}
}

func (r *RPCRequest) Encode() ([]byte, error) {
	return msgpack.Marshal(r)
}

func (r *RPCRequest) Decode(bytes []byte) error {
	return msgpack.Unmarshal(bytes, r)
}

type RPCResponse struct {
	Body []interface{} `json:"body"`
}

func NewRPCResponse(body []interface{}) *RPCResponse {
	return &RPCResponse{
		Body: body,
	}
}

func (r *RPCResponse) Encode() ([]byte, error) {
	return msgpack.Marshal(r)
}

func (r *RPCResponse) Decode(bytes []byte) error {
	return msgpack.Unmarshal(bytes, r)
}
