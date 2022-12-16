package main

import (
	"fmt"
	"log"
	"server/iface"
	"server/snet"
	"trasys/apis"
	"trasys/core"

	"net/http"
	_ "net/http/pprof"
)

func PrintAll() {
	fmt.Println("====================Stocks info======================")
	for _, tmp := range core.GlobalStocks.GetStocks() {
		fmt.Println("exchange type: ", tmp.GetExchangeType(), "; stock code: ",
			tmp.GetStockCode(), "; last price: ", tmp.GetLastPrice())
	}
	fmt.Println("====================Stocks info======================")
	fmt.Println("====================Holding info======================")
	for _, tmp := range core.GlobalHoldingRecord.GetHodingRecord() {
		s := tmp.GetStock()
		fmt.Println("client id: ", tmp.GetCliendId(), "exchange type: ", s.GetExchangeType(), "; stock code: ",
			s.GetStockCode(), "; last price: ", s.GetLastPrice(), "; holding amount: ",
			tmp.GetHoldAmount(), "; market value: ", tmp.GetMarketValue())
	}
	fmt.Println("====================Holding info======================")
}

// 当前客户端建立连接之后的hook函数
func OnConnectionAdd(conn iface.IConnection) {
	//创建一个User对象
	user := core.NewUser(conn)

	//将该连接绑定一个uid客户端的ID的属性
	conn.SetProperty("uid", user.GetId())

	//将user对象添加到用户管理模块中usermgr
	core.GlobalUserMgr.AddUser(*user)

	//给客户端发送MsgID：0的消息
	user.SyncUid()
}

// go run main.go -cpuprofile cpu.prof
// go tool pprof -http=localhost:8999 http://localhost:8999/debug/pprof/profile
// go tool pprof http://localhost:8999/debug/pprof/profile
func main() {

	//性能测试
	go func() {
		log.Println(http.ListenAndServe(":8998", nil))
	}()
	/* var cpuprofile = flag.String("cpuprofile", "", "请输入 -cpuprofile 指定cpu性能分析文件名称")

	flag.Parse()
	f, err := os.Create(*cpuprofile)
	if err != nil {
		log.Fatal("could not create CPU profile: ", err)
	}
	// StartCPUProfile为当前进程开启CPU profile。
	if err := pprof.StartCPUProfile(f); err != nil {
		log.Fatal("could not start CPU profile: ", err)
	}
	// StopCPUProfile会停止当前的CPU profile（如果有）
	defer pprof.StopCPUProfile() */
	//测试查看一下当前已经有的持仓记录和股票池
	// PrintAll()

	//创建server句柄
	s := snet.NewServer("Trading System")

	//连接创建和销毁的HOOK钩子函数
	s.SetOnConnStart(OnConnectionAdd)

	//注册一些路由业务
	s.AddRouter(1, &apis.PullHoldingRecords{}) //请求用户持仓记录
	s.AddRouter(2, &apis.Buying{})             //买入股票的请求
	s.AddRouter(3, &apis.Selling{})            //卖出股票的请求
	s.AddRouter(4, &apis.PullHistoryRecords{}) //拉取用户历史交易记录

	//启动服务
	s.Serve()
}
