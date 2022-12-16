package snet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"server/iface"
	"server/utils"
	"sync"
)

type Connection struct {
	//当前Conn隶属于哪个server
	TcpServer iface.IServer

	//当前连接的套接字
	Conn *net.TCPConn

	//当前连接的ID
	ConnID uint32

	//当前连接的状态
	isClosed bool

	// //当前连接所绑定的处理业务的方法API
	// handleAPI iface.HandleFunc

	//告知当前连接已经退出/停止  channel
	ExitChan chan bool

	//无缓冲的管道，用于读写Goroutine之间的消息通信
	msgChan chan []byte

	// //该连接处理的方法Router
	// Router iface.IRouter

	//当前的server的消息管理模块，用来绑定MsgID和对应的处理业务API关系
	MsgHandler iface.IMsgHandle

	//连接属性集合
	property map[string]interface{}

	//保护连接属性的锁
	propertyLock sync.RWMutex
}

// 连接的读业务方法
func (c *Connection) StartReader() {
	fmt.Println("[Reader Goroutine is running]")
	defer fmt.Println("connID = ", c.ConnID, " Reader is exit, remote addr is ", c.RemoteAddr().String())
	defer c.Stop()

	for {
		//读取客户端的数据到buf中，最大512字节
		// buf := make([]byte, utils.GlobalObject.MaxPackageSize)
		// _, err := c.Conn.Read(buf)
		// if err != nil && err != io.EOF {
		// 	fmt.Println("recv buf err", err)
		// 	continue
		// }
		// if err == io.EOF {
		// 	fmt.Println("this is the end", err)
		// 	break
		// }

		//调用当前连接所绑定的HandleAPI
		// if err := c.handleAPI(c.Conn, buf, cnt); err != nil {
		// 	fmt.Println("connID", c.ConnID, "handle is error", err)
		// 	break
		// }

		//创建拆包和封包的对象
		dp := NewDataPack()

		//读取客户端的Msg header二进制流 8字节
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("read msg head error", err)
			return
		}

		//拆包，得到msgID和msgDatalen 放在消息中
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack error: ", err)
			return
		}

		//根据datalen，再次读取data，放在msg.Data中
		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("read msg data error: ", err)
				return
			}
		}
		msg.SetData(data)

		//得到当前conn数据的Request请求数据
		req := Request{
			conn: c,
			msg:  msg,
		}

		// //执行注册的方法
		// go func(request iface.IRequest) {
		// 	c.Router.PreHandle(request)

		// 	c.Router.Handle(request)

		// 	c.Router.PostHandle(request)
		// }(&req)

		if utils.GlobalObject.WorkerPoolSize > 0 {
			c.MsgHandler.SendMsgToTaskQueue(&req)
		} else {
			//从路由中找到注册绑定的Conn对应的Router调用
			//根据绑定好的MsgId找到对应的api业务 执行
			go c.MsgHandler.DoMsgHandler(&req)
		}
	}
}

// 写消息的Goroutine，专门发送给客户端消息的模块
func (c *Connection) StartWriter() {
	fmt.Println("[Writer Goroutine is running]")
	defer fmt.Println(c.RemoteAddr().String(), " [conn Writer exit!]")

	//不断的阻塞等待channel的消息，进行写给客户端
	for {
		select {
		case data := <-c.msgChan:
			//有数据要写给客户端
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send data error, ", err)
				return
			}
		case <-c.ExitChan:
			//表示Reader已经退出，此时Writer也要退出
			return
		}
	}
}

// 启动连接，让当前的连接准备开始工作
func (c *Connection) Start() {
	fmt.Println("Conn Start() ... ConnID = ", c.ConnID)
	//启动从当前连接的读数据的业务
	go c.StartReader()

	//启动从当前连接写数据的业务
	go c.StartWriter()

	//按照开发者传递进来的 创建连接之后需要调用的处理业务，执行对应的Hook函数
	c.TcpServer.CallOnConnStart(c)
}

// 停止连接，结束当前连接的工作
func (c *Connection) Stop() {
	fmt.Println("Conn Stop() ... ConnID = ", c.ConnID)

	//如果当前的连接已经关闭
	if c.isClosed == true {
		return
	}
	c.isClosed = true

	//按照开发者传递进来的 销毁连接之后需要调用的处理业务，执行对应的Hook函数
	c.TcpServer.CallOnConnStop(c)

	//关闭socket连接
	c.Conn.Close()

	//告知Writer关闭
	c.ExitChan <- true

	//将当前连接从ConnMgr中摘除掉
	c.TcpServer.GetConnMgr().Remove(c)

	//回收资源
	close(c.ExitChan)
	close(c.msgChan)
}

// 获取当前连接的绑定的socket conn
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

// 获取当前连接模块的ID
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

// 获取远程客户端的TCP状态 IP port
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// 发送数据，将数据发送给远程的客户端
// 提供一个SendMsg方法，将我们要发送给客户端的数据，先进行封包，再发送
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("Connection closed when send msg")
	}

	//将data进行封包 MsgDataLen / MsgID / Data
	dp := NewDataPack()

	binaryMsg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("Pack error msg id = ", msgId)
		return errors.New("Pack error msg")
	}

	//将数据发送给客户端
	// if _, err := c.Conn.Write(binaryMsg); err != nil {
	// 	fmt.Println("Write msg id ", msgId, " error : ", err)
	// 	return errors.New("conn write error")
	// }
	c.msgChan <- binaryMsg

	return nil
}

// 设置连接属性
func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	c.property[key] = value
}

// 获取连接属性
func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	if value, ok := c.property[key]; ok {
		return value, nil
	} else {
		return nil, errors.New("no property found")
	}
}

// 移除连接属性
func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	delete(c.property, key)
}

// 初始化连接模块的方法
func NewConnection(server iface.IServer, conn *net.TCPConn, connID uint32, msgHandler iface.IMsgHandle) *Connection {
	c := &Connection{
		TcpServer:  server,
		Conn:       conn,
		ConnID:     connID,
		MsgHandler: msgHandler,
		isClosed:   false,
		msgChan:    make(chan []byte),
		ExitChan:   make(chan bool, 1),
		property:   make(map[string]interface{}),
	}

	//将conn加入到ConnManager中
	c.TcpServer.GetConnMgr().Add(c)

	return c
}
