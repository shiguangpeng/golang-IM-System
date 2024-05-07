package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	flag       int
}

// 新建一个客户端
func NewClient(serverIp string, serverPort int) *Client {
	// 创建客户端对象
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		flag:       999,
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

// DealResponse v0.9 单独开一个go routine来处理服务器回应的消息，直接显示到标准输出即可
func (this *Client) DealResponse() {
	// 阻塞方法，一旦conn中有数据，就往标准输出中写
	_, err := io.Copy(os.Stdout, this.conn)
	if err != nil {
		return
	}
}

// UpdateName v0.9 更新用户名
func (client *Client) UpdateName() bool {
	fmt.Println(">>>>>>请输入用户名：")
	fmt.Scanln(&client.Name)

	sendMsg := "rename|" + client.Name + "\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write err", err)
		return false
	}

	return true
}

// v0.9-菜单显示
func (this *Client) menu() bool {
	var flag int
	fmt.Println("1. 公聊模式")
	fmt.Println("2. 私聊模式")
	fmt.Println("3. 更新用户名")
	fmt.Println("0. 退出")

	_, err := fmt.Scanln(&flag)
	if err != nil {
		return false
	}

	if flag >= 0 && flag <= 3 {
		this.flag = flag
		return true
	} else {
		fmt.Println(">>>>> 请输入合法范围内的数字。 <<<<<<<")
		return false
	}
}

// Run 根据输入执行对应的业务
func (this *Client) Run() {
	for this.flag != 0 {
		// 直到输入正确才不循环
		for this.menu() != true {
		}
		// 根据不同的flag处理不同的业务
		switch this.flag {
		case 1:
			// 公聊模式
			break
		case 2:
			// 私聊模式
			break
		case 3:
			// 更新用户名
			this.UpdateName()
			break
		}
	}
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

	go client.DealResponse()

	fmt.Println(">>>>>> 连接服务器成功 <<<<<<<<<<")

	// 启动客户端的业务
	client.Run()

}
