package infra

import (
	"bufio"
	bytes2 "bytes"
	"encoding/binary"
	"github.com/vmihailenco/msgpack"
)

const RPCRequestBufferSize = 4096
const RPCResponseBufferSize = 4096

type Request interface {
	MethodName() string
	Params() []interface{}
	Encode() ([]byte, error)
	Decode(*bufio.Reader) error
}

type Response interface {
	Body() []interface{}
	Encode() ([]byte, error)
	Decode(*bufio.Reader) error
}

type RPCCodec struct{}

func (c RPCCodec) RPCCodecDecode(reader *bufio.Reader, r interface{}) error {
	// pick first 4 byte data which means length
	bufLength, _ := reader.Peek(4)
	bufLengthBuff := bytes2.NewBuffer(bufLength)
	var length int32
	err := binary.Read(bufLengthBuff, binary.LittleEndian, &length)
	if err != nil {
		return err
	}

	// buffer size which can read
	if int32(reader.Buffered()) < length+4 {
		return err
	}

	bytes := make([]byte, int(4+length))
	_, err = reader.Read(bytes)
	if err != nil {
		return err
	}
	return msgpack.Unmarshal(bytes[4:], r)
}

func (c RPCCodec) RPCCodecEncode(r interface{}) ([]byte, error) {
	bytes, err := msgpack.Marshal(r)
	if err != nil {
		return nil, err
	}
	length := int32(len(bytes))
	pkg := new(bytes2.Buffer)

	err = binary.Write(pkg, binary.LittleEndian, length)
	if err != nil {
		return nil, err
	}

	err = binary.Write(pkg, binary.LittleEndian, bytes)
	return pkg.Bytes(), err
}

type RPCRequest struct {
	RPCCodec
	RPCMethodName string
	RPCParams     []interface{}
}

func NewRPCRequest(methodName string, params []interface{}) *RPCRequest {
	return &RPCRequest{
		RPCMethodName: methodName,
		RPCParams:     params,
	}
}

func (r RPCRequest) MethodName() string {
	return r.RPCMethodName
}

func (r RPCRequest) Params() []interface{} {
	return r.RPCParams
}

func (r *RPCRequest) Encode() ([]byte, error) {
	return r.RPCCodecEncode(r)
}

func (r *RPCRequest) Decode(reader *bufio.Reader) error {
	return r.RPCCodecDecode(reader, r)
}

type RPCResponse struct {
	RPCCodec
	RPCBody []interface{}
}

func NewRPCResponse(body []interface{}) *RPCResponse {
	return &RPCResponse{
		RPCBody: body,
	}
}

func (r RPCResponse) Body() []interface{} {
	return r.RPCBody
}

func (r *RPCResponse) Encode() ([]byte, error) {
	return r.RPCCodecEncode(r)
}

func (r *RPCResponse) Decode(reader *bufio.Reader) error {
	return r.RPCCodecDecode(reader, r)
}
