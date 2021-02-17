package main

import (
	"fmt"
	MyClient "github.com/hejiadong/myrpc/socket/client"
)

func main() {
	client := MyClient.NewMyClient("tcp", "127.0.0.1:9999")
	b, err := client.Call("add", 0)
	if err != nil {
		fmt.Printf("err: %v", err)
	}
	fmt.Printf("result: %v", b)
	return
}
