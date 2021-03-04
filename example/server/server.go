package main

import (
	socket "github.com/hejiadong/myrpc/socket/myError"
	server "github.com/hejiadong/myrpc/socket/server"
)

type Test struct {
	C int
}

type Params struct {
	A    int
	B    int
	Test Test
}

type AddService struct {
}

func (s AddService) Add(params Params) (int, socket.MyError) {
	return params.A + params.B, nil
}

func (s AddService) Name() string {
	return "AddService"
}

func main() {
	service := AddService{}
	server := server.NewMyServer("tcp", "127.0.0.1:9999")
	server.Register(&service)
	server.Listen()
}
