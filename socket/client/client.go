//+build linux amd64

package client

import (
	"bufio"
	"fmt"
	"github.com/hejiadong/myrpc/socket/infra"
	"github.com/hejiadong/myrpc/socket/service"
	"github.com/mitchellh/mapstructure"
	"net"
	"reflect"
)

type MyClient struct {
	connect     net.Conn
	connType    string
	address     string
	name2result map[string][]reflect.Type
}

func (c MyClient) send(request infra.Request) error {
	bytes, err := request.Encode()
	if err != nil {
		return err
	}
	_, err = c.connect.Write(bytes)
	if err != nil {
		return err
	}
	return nil
}

func (c MyClient) get() (infra.Response, error) {
	reader := bufio.NewReader(c.connect)
	var response infra.RPCResponse
	err := response.Decode(reader)
	return &response, err
}

func (c MyClient) call(method string, params []interface{}) ([]interface{}, error) {
	request := infra.NewRPCRequest(method, params)
	err := c.send(request)
	if err != nil {
		return nil, err
	}
	response, err := c.get()
	if err != nil {
		return nil, err
	}
	return response.Body(), nil
}

func (c MyClient) convertParams(values []reflect.Value) ([]interface{}, error) {
	var paramInterfaces []interface{}

	start := 0
	end := len(values)

	if end-start <= 0 {
		paramInterfaces = []interface{}{}
	} else {
		paramInterfaces = make([]interface{}, end-start)
		index := 0
		for i := start; i < end; i++ {
			paramInterfaces[index] = values[i].Interface()
			index++
		}
	}
	return paramInterfaces, nil
}

func (c MyClient) convertResults(methodName string, resultInterfaces []interface{}) ([]reflect.Value, error) {
	result := make([]reflect.Value, 0)
	if len(resultInterfaces) != len(c.name2result[methodName]) {
		return result, fmt.Errorf("[MyClient]different result num betwwen remote and local")
	}
	for i := 0; i < len(resultInterfaces); i++ {
		var tmp reflect.Value
		if c.name2result[methodName][i].Kind() == reflect.Struct {
			inter := reflect.New(c.name2result[methodName][i]).Interface()
			mapstructure.Decode(resultInterfaces[i], &inter)
			tmp = reflect.ValueOf(inter).Elem()
		} else {
			tmp = reflect.ValueOf(resultInterfaces[i]).Convert(c.name2result[methodName][i])
		}
		result = append(result, tmp)
	}
	return result, nil
}

func (c MyClient) makeCallFunc(methodName string) func([]reflect.Value) []reflect.Value {
	return func(params []reflect.Value) []reflect.Value {
		paramInterfaces, _ := c.convertParams(params)

		resultInterfaces, err := c.call(methodName, paramInterfaces)
		if err != nil {
			panic(err)
		}

		result, err := c.convertResults(methodName, resultInterfaces)
		if err != nil {
			panic(err)
		}
		return result
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
			// different between t.Type and v.Type()
			v.Set(reflect.MakeFunc(t.Type, c.makeCallFunc(t.Name)))
			vt := v.Type()
			result := make([]reflect.Type, 0)
			for i := 0; i < vt.NumOut(); i++ {
				result = append(result, vt.Out(i))
			}
			c.name2result[t.Name] = result
		}
	}
}

func NewMyClient(conType string, address string) *MyClient {
	conn, err := net.Dial(conType, address)
	name2result := make(map[string][]reflect.Type)
	if err != nil {
		fmt.Printf("[NewMyClient] build conn error: %v", err)
	}
	// defer means?
	return &MyClient{
		connect:     conn,
		connType:    conType,
		address:     address,
		name2result: name2result,
	}
}
