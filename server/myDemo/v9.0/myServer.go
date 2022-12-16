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
// func (this *PingRouter) PreHandle(request iface.IRequest) {
// 	fmt.Println("Call Router PreHandle...\n")
// }

// Test Handle
func (this *PingRouter) Handle(request iface.IRequest) {
	fmt.Println("Call PingRouter Handle...")
	//先读取客户端的数据，再回写
	fmt.Println("recv from client: msgID = ", request.GetMsgID(), ", data = ", string(request.GetData()))

	//回写
	err := request.GetConnection().SendMsg(200, []byte("ping...ping...ping..."))
	if err != nil {
		fmt.Println(err)
	}
}

// Test PostHandle
// func (this *PingRouter) PostHandle(request iface.IRequest) {
// 	fmt.Println("Call Router PostHandle...")
// }

// hello test自定义路由
type HelloRouter struct {
	snet.BaseRouter
}

// Test Handle
func (this *HelloRouter) Handle(request iface.IRequest) {
	fmt.Println("Call HelloRouter Handle...")
	//先读取客户端的数据，再回写
	fmt.Println("recv from client: msgID = ", request.GetMsgID(), ", data = ", string(request.GetData()))

	//回写
	err := request.GetConnection().SendMsg(201, []byte("Hello...Hello...Hello..."))
	if err != nil {
		fmt.Println(err)
	}
}

// 创建连接之后的Hook函数
func DoConnectionBegin(conn iface.IConnection) {
	fmt.Println("===> DoConnection is Called ... ")
	if err := conn.SendMsg(202, []byte("DoConnection BEGIN")); err != nil {
		fmt.Println(err)
	}
}

// 销毁连接之后的Hook函数
func DoConnectionLost(conn iface.IConnection) {
	fmt.Println("===> DoConnection is Called ... ")
	fmt.Println("conn ID = ", conn.GetConnID(), " is Lost...")
}

func main() {
	//1 创建一个server句柄，使用server的端口
	s := snet.NewServer("version:8.0")

	//注册连接的Hook钩子函数
	s.SetOnConnStart(DoConnectionBegin)
	s.SetOnConnStop(DoConnectionLost)

	//添加自定义的Router
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloRouter{})

	//2. 启动server
	s.Serve()
}
