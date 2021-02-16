package main

import (
	"encoding/json"
	"fmt"
	"net"
)

type MyClient struct {
	connect  net.Conn
	connType string
	address  string
}

type DataFrame struct {
	Method string
	Data   interface{}
}

func (c MyClient) add(a int, b int) (interface{}, error) {
	type Params struct {
		A int
		B int
	}
	bytes, err := c.encode("add", Params{A: a, B: b})
	if err != nil {
		return 0, err
	}
	err = c.send(bytes)
	if err != nil {
		return 0, err
	}
	resultBytes, err := c.get()
	if err != nil {
		return 0, err
	}
	result, err := c.decode(resultBytes)
	if err != nil {
		return 0, err
	}
	return result, nil
}

func (c MyClient) encode(method string, data interface{}) ([]byte, error) {
	dataFrame := DataFrame{
		Method: method,
		Data:   data,
	}
	return json.Marshal(dataFrame)
}

func (c MyClient) decode(bytes []byte) (interface{}, error) {
	var out interface{}
	err := json.Unmarshal(bytes, &out)
	return out, err
}

func (c MyClient) send(bytes []byte) error {
	n, err := c.connect.Write(bytes)
	if err != nil {
		return err
	}
	println(n)
	return nil
}

func (c MyClient) get() ([]byte, error) {
	var buf [512]byte
	n, err := c.connect.Read(buf[:])
	if err != nil {
		return nil, err
	}
	return buf[:n], nil
}

func newMyClient(conType string, address string) *MyClient {
	conn, err := net.Dial(conType, address)
	if err != nil {
		fmt.Printf("[newMyClient] build conn error: %v", err)
	}
	// defer means?
	return &MyClient{
		connect:  conn,
		connType: conType,
		address:  address,
	}
}

func main() {
	client := newMyClient("tcp", "127.0.0.1:9999")
	b, err := client.add(1, 2)
	if err != nil {
		fmt.Printf("err: %v", err)
	}
	fmt.Printf("result: %v", b)
	return
}
