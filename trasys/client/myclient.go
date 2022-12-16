package main

import (
	"fmt"
	"io"
	"net"

	"server/snet"
	"trasys/pb"

	"github.com/golang/protobuf/proto"
)

// 持仓信息
type MyHoldingInfo struct {
	exchangeType string
	stockCode    string
	lastPrice    float64
	holdAmount   int32
	marketValue  float64
}

// 持仓列表信息
var MyHoldingRecords []*MyHoldingInfo

// 定义两个变量，分别是用户存储客户端的Uid和客户号的
var Cid string //客户号，这是提示用户，然用户自己输入的
var Uid int32  //用户好，这是由服务端生成，然后进行同步的

// 打印用户持仓信息
func PrintRecords() {
	fmt.Println("===================================== Your Info ========================================")
	fmt.Println("**  Your Client Id is : ", Cid)
	fmt.Println("**  Your User Id is : ", Uid)
	fmt.Println("==================================== Your Holding ======================================")
	fmt.Println("**     exchangeType      stockCode      lastPrice      holdAmount      marketValue    **")

	sum := 0.00
	for _, tmp := range MyHoldingRecords {
		fmt.Println("**      ", tmp.exchangeType, "    ", tmp.stockCode, "    ", tmp.lastPrice, "    ", tmp.holdAmount, "    ", tmp.marketValue)
		sum += tmp.marketValue
	}

	fmt.Println("====================================== Your Value =======================================")
	fmt.Println("**  Your All Market Value is : ", sum)
	fmt.Println("========================================= Menu ==========================================")
	fmt.Println("**                   [Menu] You can Enter CMD Like This                                **")
	fmt.Println("**               For Buy in: buy [exchangeType] [stockCode] [number]                   **")
	fmt.Println("**             For Sell out: sell [exchangeType] [stockCode] [number]                  **")
	fmt.Println("**                       Get Your History: history                                     **")
	fmt.Println("**                       Get Your Holding: holding                                     **")
	fmt.Println("=========================================================================================")
}

// 接受消息的结构，对不同的消息类别分别调用不同的处理结构
func RecvMsg(conn net.Conn, dp *snet.DataPack) {
	//服务器应该给我们回复一个message数据， MsgID: 1 ping ping ping

	//先读取流中的head部分，得到ID和datalen
	binaryHead := make([]byte, dp.GetHeadLen())
	if _, err := io.ReadFull(conn, binaryHead); err != nil {
		fmt.Println("read head err: ", err)
	}

	//将二进制的head拆包到msg结构体中
	msgHead, err := dp.Unpack(binaryHead)
	if err != nil {
		fmt.Println("client unpack msgHead error: ", err)
	}

	//再根据datalen进行第二次读取，将data读出来
	if msgHead.GetMsgLen() > 0 {
		msg := msgHead.(*snet.Message)
		msg.Data = make([]byte, msg.GetMsgLen())

		if _, err := io.ReadFull(conn, msg.Data); err != nil {
			fmt.Println("read msg data error: ", err)
			return
		}

		// fmt.Println("---> Recv Server Msg: ID = ", msg.Id, ", len = ", msg.DataLen, ", data = ", string(msg.Data))

		switch {
		case msg.Id == 0:
			//收到的信息是同步Uid的信息
			ClientSyncUid(msg.Id, msg.DataLen, msg.Data, dp)
		case msg.Id == 200:
			//收到的信息是拉取的客户的所有持仓记录
			ClientHoldingRecords(msg.Id, msg.DataLen, msg.Data, dp)
		case msg.Id == 201:
			//收到的信息是用户的指令执行完成的结果
			ClientReadingResult(msg.Id, msg.DataLen, msg.Data, dp)
		case msg.Id == 202:
			//收到的信息是历史交易记录
			ClientHistoryRecord(msg.Id, msg.DataLen, msg.Data, dp)
		}
	}
}

// 消息的发送结构，将客户端的信息发送出去，由不同的业务结构调用
func SendMsg(conn net.Conn, binaryMsg []byte) {
	if _, err := conn.Write(binaryMsg); err != nil {
		fmt.Println("write error: ", err)
	}
}

/*********************下面是关于读取服务端消息的实现**************************************/

// 同步服务器端的Uid，处理MsgId=0的信息
func ClientSyncUid(msgId uint32, msgDataLen uint32, msgData []byte, dp *snet.DataPack) {
	//同步用户的Uid
	newData := &pb.SyncUid{}
	err := proto.Unmarshal(msgData, newData)
	if err != nil {
		fmt.Println("unmarshal err, ", err)
	}

	Uid = newData.GetUid()
}

// 解析服务器端发送过来的持仓信息，处理MsgId=200的信息
func ClientHoldingRecords(msgId uint32, msgDataLen uint32, msgData []byte, dp *snet.DataPack) {
	records := &pb.AckHoldingInfos{}
	err := proto.Unmarshal(msgData, records)
	if err != nil {
		fmt.Println("unmarshal err, ", err)
	}

	//使用records的HInfo来初始化MyHoldingRecords
	for _, tmp := range records.GetHInfo() {
		s := &MyHoldingInfo{
			exchangeType: tmp.GetExchangeType(),
			stockCode:    tmp.GetStockCode(),
			lastPrice:    tmp.GetLastPrice(),
			holdAmount:   tmp.GetHoldAmount(),
			marketValue:  tmp.GetMarketValue(),
		}

		MyHoldingRecords = append(MyHoldingRecords, s)
	}

	fmt.Println("Your commend is succ, get your holdingrecord")
}

// 解析服务器端发送过来的持仓信息，处理MsgId=201的信息
func ClientReadingResult(msgId uint32, msgDataLen uint32, msgData []byte, dp *snet.DataPack) {
	result := &pb.ChangeStock{}
	err := proto.Unmarshal(msgData, result)
	if err != nil {
		fmt.Println("unmarshal err, ", err)
	}

	//对不同的执行结果执行不同的业务
	switch {
	case result.GetRstType() == pb.ResultType_NOSTOCK:
		//没有这样的股票
		fmt.Println("Your can not trade a non exist stock")
	case result.GetRstType() == pb.ResultType_ADDSTOCK:
		//新买入一只股票
		s := &MyHoldingInfo{
			exchangeType: result.GetExchangeType(),
			stockCode:    result.GetStockCode(),
			lastPrice:    result.GetLastPrice(),
			holdAmount:   result.GetNumber(),
			marketValue:  result.GetLastPrice() * float64(result.GetNumber()),
		}
		MyHoldingRecords = append(MyHoldingRecords, s)
	case result.GetRstType() == pb.ResultType_ADDVALUE:
		//加仓
		for i, tmp := range MyHoldingRecords {
			if tmp.exchangeType == result.GetExchangeType() && tmp.stockCode == result.GetStockCode() {
				MyHoldingRecords[i].holdAmount += result.GetNumber()
				MyHoldingRecords[i].marketValue = result.GetLastPrice() * float64(MyHoldingRecords[i].holdAmount)
				break
			}
		}
	case result.GetRstType() == pb.ResultType_DECVALUE:
		//减仓
		for i, tmp := range MyHoldingRecords {
			if tmp.exchangeType == result.GetExchangeType() && tmp.stockCode == result.GetStockCode() {
				MyHoldingRecords[i].holdAmount -= result.GetNumber()
				MyHoldingRecords[i].marketValue = result.GetLastPrice() * float64(MyHoldingRecords[i].holdAmount)
				break
			}
		}
	case result.GetRstType() == pb.ResultType_DELSTOCK:
		//抛出
		j := 0
		for _, v := range MyHoldingRecords {
			if v.exchangeType != result.GetExchangeType() || v.stockCode != result.GetStockCode() {
				MyHoldingRecords[j] = v
				j++
			}
		}
		MyHoldingRecords = MyHoldingRecords[:j]
	case result.GetRstType() == pb.ResultType_NORECORD:
		//不持有这只股票，无法进行出售
		fmt.Println("Your haven't this stock")
	case result.GetRstType() == pb.ResultType_NOEHOUGH:
		//没有足够的数量售出，你的值太大
		fmt.Println("Your haven't enough Amount to sell")
	}
}

// 解析历史交易记录
func ClientHistoryRecord(msgId uint32, msgDataLen uint32, msgData []byte, dp *snet.DataPack) {
	records := &pb.HistoryRecord{}
	err := proto.Unmarshal(msgData, records)
	if err != nil {
		fmt.Println("unmarshal err, ", err)
	}

	fmt.Println("xxxxxxxxxxxxxxxxxxxxxxxxxxxxx Your History xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx")
	fmt.Println(records.GetHisRcd())
	fmt.Println("xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx")
}

/*********************下面是关于构建发送给服务端消息的实现**************************************/

// 拉取用户持仓列表的结构，发送MsgId=1的信息
func PullList(conn net.Conn, dp *snet.DataPack) {
	//构建MsgId：1的proto结构
	proto_msg := &pb.PullHoldingInfos{
		Cid: Cid,
	}

	data, err := proto.Marshal(proto_msg)
	if err != nil {
		fmt.Println("marshal err, ", err)
	}

	//构建一个MsgId：1的结构
	binaryMsg, err := dp.Pack(snet.NewMsgPackage(1, data))
	if err != nil {
		fmt.Println("Pack error: ", err)
	}

	//发送给客户端
	SendMsg(conn, binaryMsg)
}

// 客户端执行买入的操作，发送MsgId=2的信息进行发送
func BuyIn(conn net.Conn, dp *snet.DataPack, ET string, SC string, num int32) {
	//构建MsgId：2的proto结构
	proto_msg := &pb.TradingInfo{
		ClientId:     Cid,
		ExchangeType: ET,
		StockCode:    SC,
		TraType:      pb.TradingType_BUYIN,
		Number:       num,
	}

	//序列化
	data, err := proto.Marshal(proto_msg)
	if err != nil {
		fmt.Println("marshal err, ", err)
	}

	//构建一个MsgId：2的结构
	binaryMsg, err := dp.Pack(snet.NewMsgPackage(2, data))
	if err != nil {
		fmt.Println("Pack error: ", err)
	}

	//发送给客户端
	SendMsg(conn, binaryMsg)
}

// 客户端执行卖出的操作，发送MsgId=3的信息进行发送
func SellOut(conn net.Conn, dp *snet.DataPack, ET string, SC string, num int32) {
	//构建MsgId：3的proto结构
	proto_msg := &pb.TradingInfo{
		ClientId:     Cid,
		ExchangeType: ET,
		StockCode:    SC,
		TraType:      pb.TradingType_SELLOUT,
		Number:       num,
	}

	//序列化
	data, err := proto.Marshal(proto_msg)
	if err != nil {
		fmt.Println("marshal err, ", err)
	}

	//构建一个MsgId：2的结构
	binaryMsg, err := dp.Pack(snet.NewMsgPackage(3, data))
	if err != nil {
		fmt.Println("Pack error: ", err)
	}

	//发送给客户端
	SendMsg(conn, binaryMsg)
}

// 查询历史交易记录，MsgId = 4
func GetHistory(conn net.Conn, dp *snet.DataPack) {
	//构建MsgId：4的proto结构
	proto_msg := &pb.PullHistory{
		Cid: Cid,
	}

	//序列化
	data, err := proto.Marshal(proto_msg)
	if err != nil {
		fmt.Println("marshal err, ", err)
	}

	//构建一个MsgId：4的结构
	binaryMsg, err := dp.Pack(snet.NewMsgPackage(4, data))
	if err != nil {
		fmt.Println("Pack error: ", err)
	}

	//发送给客户端
	SendMsg(conn, binaryMsg)
}

func main() {
	fmt.Println("Welcome to Stock Trading System!!!")

	//0. 直接连接远程服务器，得到一个conn连接
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("client start err, exit")
		return
	}

	//用于封包和拆包的模块
	dp := snet.NewDataPack()

	//1. 设置用户的客户号
	fmt.Print("Please Enter Your Client ID: ")
	fmt.Scanln(&Cid)

	//2. 同步用户的Uid
	fmt.Println("Sync Your Uid with Server, Please waiting....")
	RecvMsg(conn, dp)
	fmt.Println("[Uid] Your Uid is: ", Uid)

	//3. 拉取用户的持仓列表
	// 3.1. 向服务器发送拉取列表的请求
	fmt.Println("Pulling Your Record with Server, Please waiting....")
	PullList(conn, dp)

	// 3.2. 解析服务器响应的列表信息
	RecvMsg(conn, dp)

	PrintRecords()

	// 4. 用户输入命令进行操作
	for {
		//获取用户指令
		fmt.Printf(">>> ")
		var inOrOut string //交易类型，买入或者卖出
		var ET string      //股票市场
		var SC string      //股票代码
		var num int32      //数量
		fmt.Scanf("%s", &inOrOut)
		if inOrOut != "history" && inOrOut != "holding" {
			fmt.Scanf("%s %s %d", &ET, &SC, &num)
		}

		//执行相应的指令，向服务器端发送指令
		switch {
		case inOrOut == "buy":
			BuyIn(conn, dp, ET, SC, num)
		case inOrOut == "sell":
			SellOut(conn, dp, ET, SC, num)
		case inOrOut == "history":
			GetHistory(conn, dp)
		case inOrOut == "holding":
			PrintRecords()
		}

		if inOrOut != "holding" {
			//解析服务器端相应的结果
			fmt.Println("Your Commend is excuting, Please waiting...")
			RecvMsg(conn, dp)
			fmt.Println("Your Commend have been excuted, succ...")
		}
	}
}
