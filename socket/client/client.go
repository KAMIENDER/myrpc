//+build linux amd64

package client

import (
	"bufio"
	"fmt"
	"github.com/hejiadong/myrpc/socket/MyCall"
	"github.com/hejiadong/myrpc/socket/infra"
	"github.com/hejiadong/myrpc/socket/service"
	ServicePackage "github.com/hejiadong/myrpc/socket/service"
	"github.com/mitchellh/mapstructure"
	"net"
	"reflect"
	"time"
)

type MyClient struct {
	connect          net.Conn
	connType         string
	address          string
	ConnectTimeout   time.Duration
	name2serviceInfo map[string]*ServicePackage.ServiceInfo
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

func (c MyClient) call(service string, method string, params []interface{}) ([]interface{}, error) {
	request := infra.NewRPCRequest(service, method, params)
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

func (c MyClient) convertParams(serviceName string, values []reflect.Value) ([]interface{}, error) {
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

func (c MyClient) convertResults(serviceName string, methodName string, resultInterfaces []interface{}) ([]reflect.Value, error) {
	resultTypes, ok := c.name2serviceInfo[serviceName].ResultTypes(methodName)
	if !ok {
		return nil, fmt.Errorf("[MyClient]convertResults error: result not exit")
	}
	result := make([]reflect.Value, 0)
	if len(resultInterfaces) != len(resultTypes) {
		return result, fmt.Errorf("[MyClient]different result num betwwen remote and local")
	}
	for i := 0; i < len(resultInterfaces); i++ {
		var tmp reflect.Value
		if resultTypes[i].Kind() == reflect.Struct {
			inter := reflect.New(resultTypes[i]).Interface()
			mapstructure.Decode(resultInterfaces[i], &inter)
			tmp = reflect.ValueOf(inter).Elem()
		} else {
			tmp = reflect.ValueOf(resultInterfaces[i]).Convert(resultTypes[i])
		}
		result = append(result, tmp)
	}
	return result, nil
}

func (c MyClient) makeCallFunc(methodName string, call MyCall.RPCCall, serviceName string) func([]reflect.Value) []reflect.Value {
	return func(params []reflect.Value) []reflect.Value {
		paramInterfaces, _ := c.convertParams(serviceName, params)

		resultInterfaces, err := c.call(serviceName, methodName, paramInterfaces)
		if err != nil {
			panic(err)
		}

		result, err := c.convertResults(serviceName, methodName, resultInterfaces)
		if err != nil {
			panic(err)
		}
		if call != nil {
			call.SetResult(result, nil)
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
			v.Set(reflect.MakeFunc(t.Type, c.makeCallFunc(t.Name, nil, elemT.Name())))
		}
	}
	c.name2serviceInfo[service.Name()] = ServicePackage.NewServiceInfoByService(service)
}

func (c MyClient) AsyncCall(service service.RPCService, method string, params reflect.Value) MyCall.RPCCall {
	call := MyCall.NewMyCall(method, params)
	serviceType := reflect.TypeOf(service)
	oriFunc, found := serviceType.FieldByName(method)
	if !found {
		panic("[MyClient]AsyncCall: Field not found")
	}
	fun := reflect.MakeFunc(oriFunc.Type, c.makeCallFunc(method, call, serviceType.Name()))
	nowParams := make([]reflect.Value, 1)
	nowParams[0] = params
	go fun.Call(nowParams)
	return call
}

func NewMyClient(conType string, address string, timeoutSecond int64) *MyClient {
	// 建立链接是否超时
	conn, err := net.DialTimeout(conType, address, time.Duration(timeoutSecond)*time.Second)
	name2service := make(map[string]*ServicePackage.ServiceInfo)
	if err != nil {
		fmt.Printf("[NewMyClient] build conn error: %v", err)
	}
	// defer means?
	return &MyClient{
		connect:          conn,
		connType:         conType,
		address:          address,
		name2serviceInfo: name2service,
		ConnectTimeout:   time.Duration(timeoutSecond) * time.Second,
	}
}
