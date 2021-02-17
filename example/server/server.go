package main

import (
	"encoding/json"
	"fmt"
	"github.com/hejiadong/myrpc/socket/server"
)

func Add(params interface{}) (interface{}, error) {
	type out struct {
		Result int
	}
	type inType struct {
		A int
		B int
	}
	bytes, err := json.Marshal(params)
	if err != nil {
		return nil, error(fmt.Errorf("err: %v", err))
	}
	var data inType
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return nil, error(fmt.Errorf("params error"))
	}
	return out{
		Result: data.A + data.B,
	}, nil
}

func main() {
	server := socket.NewMyServer("tcp", "127.0.0.1:9999")
	server.Register(Add, "Add")
	server.Listen()
}
