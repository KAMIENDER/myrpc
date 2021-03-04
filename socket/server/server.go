package socket

import (
	"bufio"
	"fmt"
	"github.com/hejiadong/myrpc/socket/infra"
	"github.com/hejiadong/myrpc/socket/service"
	"github.com/mitchellh/mapstructure"
	"net"
	"reflect"
)

type serviceInfo struct {
	name2handler map[string]*reflect.Value
	name2params  map[string][]reflect.Type
}

func (i serviceInfo) Handler(methodName string) (*reflect.Value, bool) {
	handler, ok := i.name2handler[methodName]
	return handler, ok
}

func (i serviceInfo) ParamsTypes(methodName string) ([]reflect.Type, bool) {
	paramsTypes, ok := i.name2params[methodName]
	return paramsTypes, ok
}

func newServiceInfo(service service.RPCService) *serviceInfo {
	name2handler := make(map[string]*reflect.Value)
	name2params := make(map[string][]reflect.Type)

	serviceType := reflect.ValueOf(service)
	print("name" + serviceType.Type().Name())
	serviceType.NumMethod()

	elemV := serviceType.Elem()
	elemT := elemV.Type()
	fieldNum := elemV.NumMethod()
	for i := 0; i < fieldNum; i++ {
		t := elemT.Method(i)
		v := elemV.Method(i)
		if v.Kind() == reflect.Func {
			args := make([]reflect.Type, 0)
			for i := 0; i < v.Type().NumIn(); i++ {
				arg := v.Type().In(i)
				args = append(args, arg)
			}
			name2handler[t.Name] = &v
			name2params[t.Name] = args
		}
	}
	return &serviceInfo{
		name2handler: name2handler,
		name2params:  name2params,
	}
}

type MyServer struct {
	listener     net.Listener
	connType     string
	address      string
	name2service map[string]*serviceInfo
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

func (s *MyServer) convertParams(request infra.Request) ([]reflect.Value, error) {
	start := 0
	params := request.Params()
	end := len(params)
	paramVs := make([]reflect.Value, 0)
	funcParamsT, _ := s.name2service[request.ServiceName()].ParamsTypes(request.MethodName())
	if len(funcParamsT) != end-start {
		return nil, fmt.Errorf("[Server]dispatch error: num of args dismatch")
	}
	for i := start; i < end; i++ {
		var param reflect.Value
		if funcParamsT[i].Kind() == reflect.Struct {
			inter := reflect.New(funcParamsT[i]).Interface()
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
	service, ok := s.name2service[request.ServiceName()]
	if !ok {
		return nil, error(fmt.Errorf("[Server]dispatch error: service not exits"))
	}
	handler, ok := service.Handler(request.MethodName())
	if !ok {
		return nil, error(fmt.Errorf("[Server]dispatch error: func not exits"))
	}

	params, err := s.convertParams(request)
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
func (s *MyServer) Register(service service.RPCService) error {
	tServiceInfo := newServiceInfo(service)
	s.name2service[service.Name()] = tServiceInfo
	return nil
}

func NewMyServer(connType string, address string) *MyServer {
	listener, err := net.Listen(connType, address)
	name2service := make(map[string]*serviceInfo)
	if err != nil {
		fmt.Printf("[newMyServer]err :%v", err)
		return nil
	}
	return &MyServer{
		listener:     listener,
		connType:     connType,
		address:      address,
		name2service: name2service,
	}
}
