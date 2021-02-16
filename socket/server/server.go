package socket

import (
	"bufio"
	"fmt"
	"net"
)

type Socket interface {
	Serve(address string) error
	RegisterProcessor(fun func([]byte) error)
	process(conn net.Conn) error
	BuildConnection() (net.Conn, error)
	Send([]byte) error
	Receive() ([]byte, error)
}

type TCPSocket struct {
	listener  net.Listener
	processor func([]byte) error
	conn      net.Conn
}

func (s TCPSocket) RegisterProcessor(fun func([]byte) error) {
	s.processor = fun
}

func (s TCPSocket) process(conn net.Conn) error {
	defer conn.Close()
	for {
		reader := bufio.NewReader(conn)
		var buf [128]byte
		n, err := reader.Read(buf[:])
		if err != nil {
			fmt.Printf("[TCPSocket] process Read error: %v", err)
			return err
		}
		err = s.processor(buf[:n])
		if err != nil {
			fmt.Printf("[TCPSocket] process processor error: %v", err)
			return err
		}
	}
	return nil
}

func (s TCPSocket) Serve(address string) error {
	var err error
	s.listener, err = net.Listen("tcp", address)
	if err != nil {
		fmt.Printf("[TCPSocket] Serve listen error: %v", err)
		return err
	}
	fmt.Printf("[TCPSocket] Serve listening at %v", address)
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			fmt.Printf("[TCPSocket] Serve accept error: %v", err)
			return err
		}
		go s.process(conn)
	}
	return nil
}

func (s TCPSocket) BuildConnection(address string) (net.Conn, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Printf("[TCPSocket] BuildConnection error: %v", err)
		return nil, err
	}
	s.conn = conn
	return conn, err
}

func (s TCPSocket) Send(info []byte) error {
	_, err := s.conn.Write(info)
	if err != nil {
		fmt.Printf("[TCPSocket] Send Write error: %v", err)
	}
	return err
}

func (s TCPSocket) Receive() ([]byte, error) {
	var buf [1024]byte
	n, err := s.conn.Read(buf[:])
	if err != nil {
		fmt.Printf("[TCPSocket] Receive error: %v", err)
		return nil, err
	}
	return buf[:n], nil
}
