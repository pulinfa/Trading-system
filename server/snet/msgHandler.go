package snet

/*
* 消息管理模块的具体实现层
 */

import (
	"fmt"
	"server/iface"
	"server/utils"
	"strconv"
)

type MsgHandle struct {
	//存放每个MsgID所对应的处理方法
	Apis map[uint32]iface.IRouter

	//负责Worker任务的消息队列
	TaskQueue []chan iface.IRequest

	//工作池的数量
	WorkerPoolSize uint32
}

// 创建初始化MsgHandle方法
func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis:           make(map[uint32]iface.IRouter),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize, //从全局配置中获取，可以从配置中获取
		TaskQueue:      make([]chan iface.IRequest, utils.GlobalObject.WorkerPoolSize),
	}
}

// 调度执行对应的Router消息处理方法
func (mh *MsgHandle) DoMsgHandler(request iface.IRequest) {
	//1. 从Request中找到msgID
	handler, ok := mh.Apis[request.GetMsgID()]
	if !ok {
		fmt.Println("api msgID = ", request.GetMsgID(), "is NOT FOUND! Need register")
	}

	//根据MsgID调度对应的router业务
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

// 为消息添加具体的处理逻辑
func (mh *MsgHandle) AddRouter(msgID uint32, router iface.IRouter) {
	//1. 判断当前msg绑定的API处理方法是否已经存在
	if _, ok := mh.Apis[msgID]; ok {
		//id已经注册
		panic("repeat api, msgID = " + strconv.Itoa(int(msgID)))
	}

	//2. 添加msg和API的绑定关系
	mh.Apis[msgID] = router
	fmt.Println("Add api MsgID = ", msgID, " succ!")
}

// 启动一个worker工作池(开启工作池的动作只能发生一次，只有一个工作池)
func (mh *MsgHandle) StartWorkerPool() {
	//根据WorkerPool分别开启worker，每个worker用一个go承载
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		//一个worker启动
		//1. 当前的worker对应的channel消息队列 开辟空间 第0个worker就用第0个channel
		mh.TaskQueue[i] = make(chan iface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)

		go mh.StartOneWorker(i, mh.TaskQueue[i])
	}
}

// 启动一个Worker工作流程
func (mh *MsgHandle) StartOneWorker(workerID int, taskQueue chan iface.IRequest) {
	fmt.Println("Worker ID = ", workerID, " is started ...")

	//不断的阻塞等待对应消息队列的信息
	for {
		select {
		//如果又消息过来，出列的就是一个客户端的Request，执行当前的Request所绑定的业务
		case request := <-taskQueue:
			mh.DoMsgHandler(request)
		}
	}
}

// 将消息发送给TaskQueue，由Worker进行处理
func (mh *MsgHandle) SendMsgToTaskQueue(request iface.IRequest) {
	//将消息平均分配给不同的worker
	//根据客户端建立的ConnID来进行分配
	workerID := request.GetConnection().GetConnID() % mh.WorkerPoolSize
	fmt.Println("Add ConnID = ", request.GetConnection().GetConnID(),
		", request MsgID = ", request.GetMsgID(), " to WorkerID = ", workerID)

	//将消息发送给对应的worker的TaskQueue
	mh.TaskQueue[workerID] <- request
}
