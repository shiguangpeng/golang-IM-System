package main

func main() {
	// 请在此处编写你的代码
	server := NewServer("127.0.0.1", 8888)
	server.Start()
}
