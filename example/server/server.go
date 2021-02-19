package main

import (
	"fmt"
	"github.com/hejiadong/myrpc/socket/server"
)

func Add(a int, b int) (int, error) {
	type out struct {
		Result int64
	}
	return a + b, fmt.Errorf("test")
}

func main() {
	server := socket.NewMyServer("tcp", "127.0.0.1:9999")
	server.Register(Add, "Add")
	server.Listen()
}
