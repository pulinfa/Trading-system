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
	fmt.Println("Call Router Handle...")
	//先读取客户端的数据，再回写
	fmt.Println("recv from client: msgID = ", request.GetMsgID(), ", data = ", string(request.GetData()))

	//回写
	err := request.GetConnection().SendMsg(1, []byte("ping...ping...ping..."))
	if err != nil {
		fmt.Println(err)
	}
}

// Test PostHandle
// func (this *PingRouter) PostHandle(request iface.IRequest) {
// 	fmt.Println("Call Router PostHandle...")
// }

func main() {
	//1 创建一个server句柄，使用server的端口
	s := snet.NewServer("version:5.0")

	//添加自定义的Router
	s.AddRouter(&PingRouter{})

	//2. 启动server
	s.Serve()
}
