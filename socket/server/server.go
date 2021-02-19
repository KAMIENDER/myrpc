package socket

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/hejiadong/myrpc/socket/infra"
	"net"
	"reflect"
)

type MyServer struct {
	listener     net.Listener
	connType     string
	address      string
	name2handler map[string]*reflect.Value
	name2params  map[string][]reflect.Type
}

func (s MyServer) process(conn net.Conn) error {
	defer conn.Close()
	for {
		reader := bufio.NewReader(conn)
		var buf [512]byte
		n, err := reader.Read(buf[:])
		if err != nil {
			return err
		}
		dataFrame, err := s.decode(buf[:n])
		if err != nil {
			return err
		}
		result, err := s.dispatch(*dataFrame)
		if err != nil {
			return err
		}
		err = s.send(result, conn)
		if err != nil {
			return err
		}
	}
}

func (s MyServer) send(result []reflect.Value, conn net.Conn) error {
	body := make([]interface{}, 0)
	for i := 0; i < len(result); i++ {
		body = append(body, result[i].Interface())
	}
	response := infra.NewRPCResponse(body)
	buf, err := json.Marshal(response)
	if err != nil {
		return err
	}
	_, err = conn.Write(buf)
	return err
}

func (s MyServer) dispatch(request infra.RPCRequest) ([]reflect.Value, error) {
	handler, ok := s.name2handler[request.MethodName]
	if !ok {
		return nil, error(fmt.Errorf("[Server]dispatch error: func not exits"))
	}
	start := 0
	end := len(request.Params)
	paramVs := make([]reflect.Value, 0)
	funcParamsT := s.name2params[request.MethodName]
	if len(funcParamsT) != end-start {
		return nil, fmt.Errorf("[Server]dispatch error: num of args dismatch")
	}
	for i := start; i < end; i++ {
		param := reflect.ValueOf(request.Params[i]).Convert(funcParamsT[i])
		paramVs = append(paramVs, param)
	}
	result := handler.Call(paramVs)
	return result, nil
}

func (s MyServer) decode(bytes []byte) (*infra.RPCRequest, error) {
	var request infra.RPCRequest
	err := json.Unmarshal(bytes, &request)
	if err != nil {
		return nil, err
	}
	return &request, err
}

func (s MyServer) Listen() error {
	fmt.Printf("[MyServer]start listening at %v", s.listener.Addr())
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			return err
		}
		go s.process(conn)
	}
}
func (s MyServer) Register(handler interface{}, name string) error {
	handlerV := reflect.ValueOf(handler)
	handlerT := reflect.TypeOf(handler)
	if handlerT.Kind() != reflect.Func {
		return error(fmt.Errorf("[MyServer Register]error: Not Func type"))
	}
	args := make([]reflect.Type, 0)
	for i := 0; i < handlerT.NumIn(); i++ {
		arg := handlerT.In(i)
		args = append(args, arg)
	}
	s.name2handler[name] = &handlerV
	s.name2params[name] = args
	return nil
}

func NewMyServer(connType string, address string) *MyServer {
	listener, err := net.Listen(connType, address)
	name2handler := make(map[string]*reflect.Value)
	name2params := make(map[string][]reflect.Type)
	if err != nil {
		fmt.Printf("[newMyServer]err :%v", err)
		return nil
	}
	return &MyServer{
		listener:     listener,
		connType:     connType,
		address:      address,
		name2handler: name2handler,
		name2params:  name2params,
	}
}
