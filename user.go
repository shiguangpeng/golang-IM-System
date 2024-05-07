package main

import (
	"net"
	"strings"
)

// 用户类跟server交互
type User struct {
	Name string
	Addr string
	// string类型的管道channel
	C    chan string
	conn net.Conn
	// 用户绑定的服务器类型指针
	server *Server
}

func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name: userAddr,
		Addr: userAddr,
		C:    make(chan string),
		conn: conn,

		server: server,
	}
	// 启动监听当前user channel的go程
	go user.ListMessage()
	return user
}

// v0.4新增：用户的上线业务
func (this *User) Online() {
	this.server.mapLock.Lock()
	this.server.OnlineMap[this.Name] = this
	this.server.mapLock.Unlock()
	this.server.Broadcast(this, "上线")
}

// v0.4新增：用户下线业务
func (this *User) Offline() {
	this.server.mapLock.Lock()
	delete(this.server.OnlineMap, this.Name)
	this.server.mapLock.Unlock()
	this.server.Broadcast(this, "下线")
}

// v0.5 给指定用户发消息
func (this *User) SendMsg(msg string) {
	write, err := this.conn.Write([]byte(msg))
	if err != nil && write > 0 {
		return
	}
}

// v0.4新增：用户处理消息的业务
// v0.5新增：输入特定消息，如who，可以查询当前在线用户有哪些
func (this *User) DoMessage(msg string) {

	// v0.5
	// 查询当前用户都有哪些
	if msg == "who" {
		this.server.mapLock.Lock()
		for _, user := range this.server.OnlineMap {
			onlineMsg := "[" + user.Addr + "]" + user.Name + "上线...\n"
			this.SendMsg(onlineMsg)
		}
		this.server.mapLock.Unlock()
	} else if len(msg) > 7 && msg[0:7] == "rename" {
		//v0.6 修改用户名
		newName := strings.Split(msg, "|")[1]
		// 判断要修改的名字是否存在
		_, ok := this.server.OnlineMap[newName]
		if ok {
			this.SendMsg("用户名已经存在")
		} else {
			// V0.6
			// 更新用户名
			this.server.mapLock.Lock()
			this.server.OnlineMap[newName] = this
			this.server.mapLock.Unlock()
			this.Name = newName
			this.SendMsg("您已经更新用户名：" + this.Name + "\n")
		}
	} else if len(msg) > 4 && msg[:3] == "to|" {
		// v0.7 私聊功能
		// 消息模式：to|张三|消息内容
		// 1获取用户名
		remoteName := strings.Split(msg, "|")[1]
		if remoteName == "" {
			this.SendMsg("消息格式不正确，请使用 to|张三|消息内容。\n")
			return
		}
		// 2根据得到的用户名，得到对方User对象
		remoteUser, ok := this.server.OnlineMap[remoteName]
		if !ok {
			this.SendMsg("该用户不存在。")
			return
		}
		// 3获取消息内容
		content := strings.Split(msg, "|")[2]
		if content == "" {
			this.SendMsg("无消息内容，请重发！\n")
			return
		}
		remoteUser.SendMsg(this.Name + "对您说：" + content)
	} else {
		// v0.4
		this.server.Broadcast(this, msg)
	}
}

// 监听当前User channel的方法，一旦有消息，就发送出去
func (this *User) ListMessage() {
	for {
		msg := <-this.C
		this.conn.Write([]byte(msg + "\n"))
	}
}
