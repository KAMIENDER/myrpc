//+build linux amd64

package client

import (
	"fmt"
	"github.com/hejiadong/myrpc/socket/infra"
	"github.com/hejiadong/myrpc/socket/service"
	"github.com/mitchellh/mapstructure"
	"github.com/vmihailenco/msgpack"
	"net"
	"reflect"
)

type MyClient struct {
	connect     net.Conn
	connType    string
	address     string
	name2result map[string][]reflect.Type
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
	var buf [10000]byte
	n, err := c.connect.Read(buf[:])
	if err != nil {
		return nil, err
	}
	return buf[:n], nil
}

func (c MyClient) call(method string, params []interface{}) ([]interface{}, error) {
	request := infra.NewRPCRequest(method, params)
	bytes, err := msgpack.Marshal(&request)
	if err != nil {
		return nil, err
	}
	err = c.send(bytes)
	if err != nil {
		return nil, err
	}
	resultBytes, err := c.get()
	if err != nil {
		return nil, err
	}
	var response infra.RPCResponse
	err = msgpack.Unmarshal(resultBytes, &response)
	if err != nil {
		return nil, err
	}
	return response.Body, nil
}

func (c MyClient) makeCallFunc(methodName string) func([]reflect.Value) []reflect.Value {
	return func(in []reflect.Value) []reflect.Value {
		var paramInterface []interface{}

		start := 0
		end := len(in)

		if end-start <= 0 {
			paramInterface = []interface{}{}
		} else {
			paramInterface = make([]interface{}, end-start)
			index := 0
			for i := start; i < end; i++ {
				paramInterface[index] = in[i].Interface()
				index++
			}
		}

		result, err := c.call(methodName, paramInterface)
		if err != nil {
			panic(err)
		}

		out := make([]reflect.Value, 0)
		if len(result) != len(c.name2result[methodName]) {
			out = append(out, reflect.ValueOf(fmt.Errorf("[MyClient]different out num betwwen remote and local")))
			return out
		}
		for i := 0; i < len(result); i++ {
			var tmp reflect.Value
			if reflect.ValueOf(result[i]).Kind() == reflect.Interface {
				tmp = reflect.New(c.name2result[methodName][i])
				mapstructure.Decode(result[i], &tmp)
			} else {
				tmp = reflect.ValueOf(result[i]).Convert(c.name2result[methodName][i])
			}
			out = append(out, tmp)
		}
		return out
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
