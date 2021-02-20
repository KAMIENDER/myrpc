//+build linux darwin windows

package infra

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

type RPCResponse struct {
	Body []interface{} `json:"body"`
}

func NewRPCResponse(body []interface{}) *RPCResponse {
	return &RPCResponse{
		Body: body,
	}
}
