package client

import (
	"encoding/json"
	"fmt"
	"github.com/hejiadong/myrpc/socket/infra"
	"github.com/hejiadong/myrpc/socket/service"
	"net"
	"reflect"
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

func (c MyClient) call(method string, params interface{}) (interface{}, error) {
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

func (c MyClient) makeCallFunc(methodName string) func([]reflect.Value) []reflect.Value {
	return func(values []reflect.Value) []reflect.Value {
		result, err := c.call(methodName, values)
		fmt.Printf("%v %v", result, err)
		return []reflect.Value{
			reflect.ValueOf(result),
			reflect.ValueOf(err),
		}
	}
}

func (c MyClient) RegisterService(service service.RPCService) {
	serviceType := reflect.ValueOf(service)

	elemV := serviceType.Elem()
	elemT := elemV.Type()
	fieldNum := elemV.NumField()
	for i := 0; i < fieldNum; i++ {
		t := elemT.Field(i)
		v := elemV.Field(i)
		if v.Kind() == reflect.Func && v.CanSet() && v.IsValid() {
			v.Set(reflect.MakeFunc(t.Type, c.makeCallFunc(t.Name)))
		}
	}
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
