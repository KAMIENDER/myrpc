package main

import (
	socket "github.com/hejiadong/myrpc/socket/myError"
	server "github.com/hejiadong/myrpc/socket/server"
)

func Add(a int, b int) (int, socket.MyError) {
	return a + b, socket.NewRPCError("test")
}

func main() {
	server := server.NewMyServer("tcp", "127.0.0.1:9999")
	server.Register(Add, "Add")
	server.Listen()
}
