package main

import (
	"fmt"
	"net"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
}

// 新建一个客户端
func NewClient(serverIp string, serverPort int) *Client {
	// 创建客户端对象
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
	}
	// 连接server
	dial, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net.Dial error", err)
		return nil
	}
	client.conn = dial
	return client
}

// 连接到server

func main() {
	var client *Client = NewClient("127.0.0.1", 8888)
	if client == nil {
		fmt.Println(">>>>>> 连接服务器失败 <<<<<<<<<<")
		return
	}
	fmt.Println(">>>>>> 连接服务器成功 <<<<<<<<<<")

}
