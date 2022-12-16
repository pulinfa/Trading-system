package snet

import (
	"fmt"
	"net"

	"server/iface"
	"server/utils"
)

// IServer的接口的实现，定义一个Server的服务器模块
type Server struct {
	//定义属性，名称，ip版本，监听的ip，监听的端口
	Name      string
	IPVersion string
	IP        string
	Port      int

	// //给当前的Server添加一个router，server注册的来凝结对应的业务
	// Router iface.IRouter

	//当前的server的消息管理模块，用来绑定MsgID和对应的处理业务API关系
	MsgHandler iface.IMsgHandle

	//该server的连接管理器
	ConnMgr iface.IConnManager

	//创建连接之后自动调用的Hook函数
	OnConnStart func(conn iface.IConnection)

	//销毁连接之前自动调用的Hook函数
	OnConnStop func(conn iface.IConnection)
}

// 业务代码，定义当前客户端连接所绑定的handle api（后期应该自定义）
// func callBackToClient(conn *net.TCPConn, data []byte, cnt int) error {
// 	//这个代码主要是回显
// 	fmt.Println("[Conn Handle] CallbackToClient...")
// 	if _, err := conn.Write(data[:cnt]); err != nil {
// 		fmt.Println("wirte back buf err", err)
// 		return errors.New("CallBackToClient error")
// 	}

// 	return nil
// }

// 启动服务器
func (s *Server) Start() {
	fmt.Printf("[Start] Server Name: %s, Listenner at IP : %s, Port : %d, is starting\n",
		utils.GlobalObject.Name, utils.GlobalObject.Host, utils.GlobalObject.TcpPort)
	fmt.Printf("[Start] Version %s, MaxConn: %d, MaxPackageSize: %d\n",
		utils.GlobalObject.Version,
		utils.GlobalObject.MaxConn,
		utils.GlobalObject.MaxPackageSize)

	go func() {
		//0. 开启消息队列和worker工作池
		s.MsgHandler.StartWorkerPool()

		//1. 获取一个TCP的Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("Reslove tcp addr error!")
			return
		}

		//2. 监听服务器的地址
		listenner, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("Listen "+s.IPVersion+" err ", err)
			return
		}

		fmt.Println("Start server succe " + s.Name + " succ, Listenning...")

		var cid uint32 //记录连接的ID
		cid = 0

		//3. 阻塞的等待客户端的连接，处理客户端连接服务（读写）
		for {
			conn, err := listenner.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err", err)
				continue
			}

			//设置最大的连接数量的判断，如果超过最大连接，则关闭
			if s.ConnMgr.Len() >= utils.GlobalObject.MaxConn {
				//给客户端相应一个超出最大连接的错误包
				fmt.Println("Too many Connections")
				conn.Close()
				continue
			}

			//将处理新连接的业务方法和conn进行绑定
			dealConn := NewConnection(s, conn, cid, s.MsgHandler) //监听到一个端口，这个时候就注册一个连接，对连接进行处理
			cid++

			//启动当前的连接的业务处理
			go dealConn.Start()
		}
	}()
}

// 停止服务器
func (s *Server) Stop() {
	// 将服务器的一些资源、状态或者一些已经开辟的连接信息，进行停止和回收
	fmt.Println("[Stop] server name: ", s.Name)
	s.ConnMgr.ClearConn()
}

// 运行服务器
func (s *Server) Serve() {
	s.Start()

	stopruning := make(chan int)
	defer close(stopruning)

	go func() {
		var tmp int
		fmt.Scanf("%d", &tmp)
		fmt.Println("your ask to stop")

		if tmp == 1 {
			stopruning <- 1
			return
		}
	}()

	//阻塞状态
	select {
	case _, ok := <-stopruning:
		if ok {
			fmt.Println("get stop")
			s.Stop()
			return
		}
	}
}

// 路由功能：给当前的服务注册一个路由方法，供客户端的连接使用
func (s *Server) AddRouter(msgID uint32, router iface.IRouter) {
	s.MsgHandler.AddRouter(msgID, router)
	fmt.Println("Add Router Succ")
}

func (s *Server) GetConnMgr() iface.IConnManager {
	return s.ConnMgr
}

// 初始化Server模块的方法
func NewServer(name string) iface.IServer {
	s := &Server{
		Name:       utils.GlobalObject.Name,
		IPVersion:  "tcp4",
		IP:         utils.GlobalObject.Host,
		Port:       utils.GlobalObject.TcpPort,
		MsgHandler: NewMsgHandle(),
		ConnMgr:    NewConnManager(),
	}

	return s
}

// ​注册OnConnStart钩子函数的方法
func (s *Server) SetOnConnStart(hookFunc func(connection iface.IConnection)) {
	s.OnConnStart = hookFunc
}

// ​注册OnConnStop钩子函数的方法
func (s *Server) SetOnConnStop(hookFunc func(connection iface.IConnection)) {
	s.OnConnStop = hookFunc
}

// ​调用OnConnStop钩子函数的方法
func (s *Server) CallOnConnStart(conn iface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("---> Call OnConnStart()...")
		s.OnConnStart(conn)
	}
}

// ​调用OnConnStop钩子函数的方法
func (s *Server) CallOnConnStop(conn iface.IConnection) {
	if s.OnConnStop != nil {
		fmt.Println("---> Call OnConnStop()...")
		s.OnConnStop(conn)
	}
}
