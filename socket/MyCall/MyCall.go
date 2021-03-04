package MyCall

import (
	"reflect"
	"sync"
)

type RPCCall interface {
	Result() []reflect.Value
	SetResult([]reflect.Value, error)
	Error() error
}

type MyCall struct {
	Done bool
	ch chan bool
	method string
	err error
	Params reflect.Value
	result []reflect.Value
	mutex sync.Mutex
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

func (c *MyCall) SetResult(value []reflect.Value, err error)  {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.err = err
	c.result = value
	c.ch <- true
}

func (c MyCall) Error() error {
	if c.Done {
		return c.err
	}
	c.Done =<- c.ch
	return c.err
}

