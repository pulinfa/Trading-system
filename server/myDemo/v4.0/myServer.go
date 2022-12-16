package main

import (
	"fmt"
	"server/iface"
	"server/snet"
)

// ping test自定义路由
type PingRouter struct {
	snet.BaseRouter
}

// Test PreHandle
func (this *PingRouter) PreHandle(request iface.IRequest) {
	fmt.Println("Call Router PreHandle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("before ping...\n"))
	if err != nil {
		fmt.Println("call back before ping error")
	}
}

// Test Handle
func (this *PingRouter) Handle(request iface.IRequest) {
	fmt.Println("Call Router Handle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("ping... ping... ping...\n"))
	if err != nil {
		fmt.Println("call back ping error")
	}
}

// Test PostHandle
func (this *PingRouter) PostHandle(request iface.IRequest) {
	fmt.Println("Call Router PostHandle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("after ping...\n"))
	if err != nil {
		fmt.Println("call back after ping error")
	}
}

func main() {
	//1 创建一个server句柄，使用server的端口
	s := snet.NewServer("version:3.0")

	//添加自定义的Router
	s.AddRouter(&PingRouter{})

	//2. 启动server
	s.Serve()
}
