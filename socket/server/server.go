package socket

import (
	"bufio"
	"fmt"
	"github.com/hejiadong/myrpc/socket/infra"
	"github.com/mitchellh/mapstructure"
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

func (s *MyServer) process(con net.Conn) error {
	defer con.Close()
	for {
		request, err := s.get(con)
		if err != nil {
			return err
		}
		response, err := s.dispatch(request)
		if err != nil {
			return err
		}
		err = s.send(response, con)
		if err != nil {
			return err
		}
	}
}

func (s MyServer) get(con net.Conn) (infra.Request, error) {
	var request infra.RPCRequest

	reader := bufio.NewReader(con)

	err := request.Decode(reader)
	return &request, err
}

func (s MyServer) send(response infra.Response, conn net.Conn) error {
	buf, err := response.Encode()
	if err != nil {
		return err
	}
	_, err = conn.Write(buf)
	return err
}

func (s *MyServer) convertParams(methodName string, params []interface{}) ([]reflect.Value, error) {
	start := 0
	end := len(params)
	paramVs := make([]reflect.Value, 0)
	funcParamsT := s.name2params[methodName]
	if len(funcParamsT) != end-start {
		return nil, fmt.Errorf("[Server]dispatch error: num of args dismatch")
	}
	for i := start; i < end; i++ {
		var param reflect.Value
		if s.name2params[methodName][i].Kind() == reflect.Struct {
			inter := reflect.New(s.name2params[methodName][i]).Interface()
			mapstructure.Decode(params[i], &inter)
			param = reflect.ValueOf(inter).Elem()
		} else {
			param = reflect.ValueOf(params[i]).Convert(funcParamsT[i])
		}
		paramVs = append(paramVs, param)
	}
	return paramVs, nil
}

func (s MyServer) convertResult(result []reflect.Value) ([]interface{}, error) {
	resultInterfaces := make([]interface{}, 0)
	for i := 0; i < len(result); i++ {
		resultInterfaces = append(resultInterfaces, result[i].Interface())
	}
	return resultInterfaces, nil
}

func (s *MyServer) dispatch(request infra.Request) (infra.Response, error) {
	handler, ok := s.name2handler[request.MethodName()]
	if !ok {
		return nil, error(fmt.Errorf("[Server]dispatch error: func not exits"))
	}

	params, err := s.convertParams(request.MethodName(), request.Params())
	if err != nil {
		return nil, err
	}
	result := handler.Call(params)

	resultInterfaces, err := s.convertResult(result)
	response := infra.NewRPCResponse(request.ServiceName(), resultInterfaces)

	return response, nil
}

func (s *MyServer) Listen() error {
	fmt.Printf("[MyServer]start listening at %v", s.listener.Addr())
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			return err
		}
		go s.process(conn)
	}
}
func (s *MyServer) Register(handler interface{}, name string) error {
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
