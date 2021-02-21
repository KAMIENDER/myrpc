package main

import (
	"fmt"
	MyClient "github.com/hejiadong/myrpc/socket/client"
	socket "github.com/hejiadong/myrpc/socket/myError"
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
	Add func(params Params) (int, socket.RPCError)
}

func (s *AddService) Name() string {
	return "Add"
}

func main() {
	var s AddService
	client := MyClient.NewMyClient("tcp", "127.0.0.1:9999")
	client.RegisterService(&s)
	b, err := s.Add(Params{
		A: 1,
		B: 2,
		Test: Test{
			C: 10,
		},
	})
	fmt.Printf("%v %v", b, err)
}
