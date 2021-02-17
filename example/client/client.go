package main

import (
	"fmt"
	MyClient "github.com/hejiadong/myrpc/socket/client"
)

type AddService struct {
	Add func(a int, b int) (int, error)
}

func (s *AddService) Name() string {
	return "Add"
}

func main() {
	var s AddService
	client := MyClient.NewMyClient("tcp", "127.0.0.1:9999")
	client.RegisterService(&s)
	b, err := s.Add(1, 2)
	fmt.Printf("%v %v", b, err)
}
