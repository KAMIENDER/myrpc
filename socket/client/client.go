package client

import (
	"encoding/json"
	"fmt"
	"github.com/hejiadong/myrpc/socket/infra"
	"net"
)

type MyClient struct {
	connect  net.Conn
	connType string
	address  string
}

func (c MyClient) encode(method string, data interface{}) ([]byte, error) {
	dataFrame := infra.DataFrame{
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

func (c MyClient) Call(method string, params interface{}) (interface{}, error) {
	bytes, err := c.encode(method, params)
	if err != nil {
		return nil, err
	}
	err = c.send(bytes)
	if err != nil {
		return nil, err
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

func NewMyClient(conType string, address string) *MyClient {
	conn, err := net.Dial(conType, address)
	if err != nil {
		fmt.Printf("[NewMyClient] build conn error: %v", err)
	}
	// defer means?
	return &MyClient{
		connect:  conn,
		connType: conType,
		address:  address,
	}
}
