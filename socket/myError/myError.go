package socket

import (
	"runtime"
)

var ErrorStackBufferSize = 4096

type MyError interface {
	Error() string
	Stack() string
}

type RPCError struct {
	Info      string        `json:"info"`
	StackInfo string        `json:"stack"`
	Extra     []interface{} `json:"Extra"`
}

func NewRPCError(info string) *RPCError {
	stack := make([]byte, ErrorStackBufferSize) //4KB
	runtime.Stack(stack, true)
	return &RPCError{
		Info:      info,
		StackInfo: string(stack),
	}
}

func (e RPCError) Error() string {
	return "info: " + e.Info + "\nstack:\n" + e.StackInfo
}

func (e RPCError) Stack() string {
	return e.StackInfo
}
