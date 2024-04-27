package main

import "net"

// 用户类跟server交互
type User struct {
	Name string
	Addr string
	// string类型的管道channel
	C    chan string
	conn net.Conn
}

func NewUser(conn net.Conn) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name: userAddr,
		Addr: userAddr,
		C:    make(chan string),
		conn: conn,
	}
	// 启动监听当前user channel的go程
	go user.ListMessage()
	return user
}

// 监听当前User channel的方法，一旦有消息，就发送出去
func (this *User) ListMessage() {
	for {
		msg := <-this.C
		this.conn.Write([]byte(msg + "\n"))
	}
}
