package MyCall

import "reflect"

type RPCCall interface {
	Result() []reflect.Value

}

type MyCall struct {
	Done bool
	ch chan bool
	method string
	Params reflect.Value
	result []reflect.Value
}

func NewMyCall(method string, params reflect.Value) *MyCall {
	ch := make(chan bool, 1)
	return &MyCall{
		ch: ch,
		method: method,
		Done: false,
		Params: params,
		result: nil,
	}
}

func (c *MyCall) Result() []reflect.Value {
	if c.Done {
		return c.result
	}
	c.Done =<- c.ch
	return c.result
}

func (c *MyCall) SetResult(value []reflect.Value)  {
	c.result = value
	c.ch <- true
}

