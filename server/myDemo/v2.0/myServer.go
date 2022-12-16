package main

import "server/snet"

func main() {
	//1 创建一个server句柄，使用server的端口
	s := snet.NewServer("version:2.0")

	//2. 启动server
	s.Serve()
}
