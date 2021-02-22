//+build linux darwin windows

package infra

import "github.com/vmihailenco/msgpack"

const RPCRequestBufferSize = 4096
const RPCResponseBufferSize = 4096

type Request interface {
	MethodName() string
	Params() []interface{}
	Encode() ([]byte, error)
	Decode([]byte) error
}

type Response interface {
	Body() []interface{}
	Encode() ([]byte, error)
	Decode([]byte) error
}

type RPCRequest struct {
	RPCMethodName string        `json:"method_name"`
	RPCParams     []interface{} `json:"params"`
}

func NewRPCRequest(methodName string, params []interface{}) *RPCRequest {
	return &RPCRequest{
		RPCMethodName: methodName,
		RPCParams:     params,
	}
}

func (r RPCRequest) MethodName() string {
	return r.RPCMethodName
}

func (r RPCRequest) Params() []interface{} {
	return r.RPCParams
}

func (r *RPCRequest) Encode() ([]byte, error) {
	return msgpack.Marshal(r)
}

func (r *RPCRequest) Decode(bytes []byte) error {
	return msgpack.Unmarshal(bytes, &r)
}

type RPCResponse struct {
	RPCBody []interface{} `json:"body"`
}

func NewRPCResponse(body []interface{}) *RPCResponse {
	return &RPCResponse{
		RPCBody: body,
	}
}

func (r RPCResponse) Body() []interface{} {
	return r.RPCBody
}

func (r *RPCResponse) Encode() ([]byte, error) {
	return msgpack.Marshal(r)
}

func (r *RPCResponse) Decode(bytes []byte) error {
	return msgpack.Unmarshal(bytes, r)
}
