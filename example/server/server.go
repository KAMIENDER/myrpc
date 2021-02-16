package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
)

type MyServer struct {
	listener    net.Listener
	connType    string
	address     string
	string2func map[string]func(interface{}) (interface{}, error)
}

type DataFrame struct {
	Method string
	Data   interface{}
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
		result, err := s.dispatch(dataFrame)
		if err != nil {
			return err
		}
		err = s.send(result, conn)
		if err != nil {
			return err
		}
	}
}

func (s MyServer) send(result interface{}, conn net.Conn) error {
	buf, err := json.Marshal(result)
	if err != nil {
		return err
	}
	_, err = conn.Write(buf)
	return err
}

func (s MyServer) dispatch(dataFrame DataFrame) (interface{}, error) {
	handler, ok := s.string2func[dataFrame.Method]
	if !ok {
		return nil, error(fmt.Errorf("func not exits"))
	}
	return handler(dataFrame.Data)
}

func (s MyServer) decode(bytes []byte) (DataFrame, error) {
	var a DataFrame
	err := json.Unmarshal(bytes, &a)
	return a, err
}

func (s MyServer) listen() error {
	fmt.Printf("[MyServer]start listening at %v", s.listener.Addr())
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			return err
		}
		go s.process(conn)
	}
}
func (s MyServer) register(handler func(interface{}) (interface{}, error), name string) error {
	s.string2func[name] = handler
	return nil
}

func newMyServer(connType string, address string) *MyServer {
	listener, err := net.Listen(connType, address)
	string2func := make(map[string]func(interface{}) (interface{}, error))
	if err != nil {
		fmt.Printf("[newMyServer]err :%v", err)
		return nil
	}
	return &MyServer{
		listener:    listener,
		connType:    connType,
		address:     address,
		string2func: string2func,
	}
}

func add(params interface{}) (interface{}, error) {
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
	server := newMyServer("tcp", "127.0.0.1:9999")
	server.register(add, "add")
	server.listen()
}
