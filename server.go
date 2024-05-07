package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type Server struct {
	Ip   string
	Port int

	// 在线用户列表
	OnlineMap map[string]*User
	mapLock   sync.RWMutex

	// 消息广播
	Message chan string
}

// 创建一个Server接口，返回的是一个Server结构体指针
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return server
}

// 监听Message广播消息channel的goroutine, 一旦有消息就发送给全部的在线user
func (this *Server) ListenMessager() {
	for {
		msg := <-this.Message

		// 将msg发送给全部的在线user
		this.mapLock.Lock()
		for _, cli := range this.OnlineMap {
			cli.C <- msg
		}
		this.mapLock.Unlock()
	}
}

// 广播用户上线消息
func (this *Server) Broadcast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg
	this.Message <- sendMsg
}

// handler
func (this *Server) Handler(conn net.Conn) {
	// handler
	//fmt.Println("连接建立成功！")
	// 用户上线，将用户加入到onlinemap中
	//user := NewUser(conn)
	// 用户与server关联
	user := NewUser(conn, this)
	user.Online()

	// 监听用户是否活跃的channel
	isLive := make(chan bool)

	// 接受客户端发送的消息
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				user.Offline()
				//this.Broadcast(user, "下线")
				return
			}
			if err != nil && err != io.EOF {
				fmt.Println("Conn Read err:", err)
				return
			}

			// 去除尾零，在控制台按下enter时，发送消息带的\n
			msg := string(buf[:n-1])
			user.DoMessage(msg)
			// 广播客户端的消息
			//this.Broadcast(user, msg)

			// 用户的任意操作代表当前用户活跃
			isLive <- true
		}
	}()

	// 当前handler阻塞
	// v0.7 超时强踢功能
	// 相当于while true
	for {
		select {
		case <-isLive:
			// 当前用户是活跃的，重置定时器
			// 不做任何事，为了触发下面的定时器
		case <-time.After(time.Second * 100):
			// 已经超时
			// 将当前的User强制关闭
			user.SendMsg("你被踢了。")

			// 销毁用的资源
			close(user.C)
			// 关闭连接
			err := conn.Close()
			if err != nil {
				return
			}
		}
	}

}

// 启动服务器的接口
func (this *Server) Start() {
	// socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("net.Listen err:", err)
	}
	// close listen socket
	defer listener.Close()

	// 启动监听Message的goroutine
	go this.ListenMessager()

	// accept，死循环
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept err:", err)
			continue
		}
		// do handle, go routine
		go this.Handler(conn)
	}
}
