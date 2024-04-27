package main

import "net"

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

// v0.4新增：用户处理消息的业务
func (this *User) DoMessage(msg string) {
	this.server.Broadcast(this, msg)
}

// 监听当前User channel的方法，一旦有消息，就发送出去
func (this *User) ListMessage() {
	for {
		msg := <-this.C
		this.conn.Write([]byte(msg + "\n"))
	}
}
