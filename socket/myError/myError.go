package socket

import (
	"runtime/debug"
)

var ErrorStackBufferSize = 4096

type MyError interface {
	Error() string
	Stack() string
}

type RPCError struct {
	Info      string
	StackInfo string
	Extra     []interface{}
}

func NewRPCError(info string) *RPCError {
	stack := debug.Stack()
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
