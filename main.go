package main // import "github.com/hejiadong/myrpc"
import (
	"fmt"
	socket "github.com/hejiadong/myrpc/socket/server"
)

func test(info []byte) error {
	fmt.Printf("%v", string(info))
	return nil
}

func main() {
	var server socket.TCPSocket
	server.RegisterProcessor(test)
	server.Serve("127.0.0.1:8080")
	fmt.Println("hello")
}
