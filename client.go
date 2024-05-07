package main

import (
	"flag"
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

var serverIP string
var serverPort int

// 命令行解析
func init() {
	flag.StringVar(&serverIP, "ip", "127.0.0.1", "设置服务器IP地址，默认是127.0.0.1")
	flag.IntVar(&serverPort, "port", 8888, "设置服务器端口号地址，默认是8888")
}
func main() {
	// 解析命令行
	flag.Parse()

	var client *Client = NewClient(serverIP, serverPort)
	if client == nil {
		fmt.Println(">>>>>> 连接服务器失败 <<<<<<<<<<")
		return
	}
	fmt.Println(">>>>>> 连接服务器成功 <<<<<<<<<<")

}
