package core

import (
	"fmt"
	"server/iface"
	"sync"

	"trasys/pb"

	"github.com/golang/protobuf/proto"
)

type User struct {
	Uid  int32             //用户的ID
	Conn iface.IConnection //当前用户的连接（用于和客户端的连接）
}

// User ID的生成器
var UidGen int32 = 1  //用来产生用户ID的计数器
var IdLock sync.Mutex //保护UidGen的Mutex

func NewUser(conn iface.IConnection) *User {
	//生成一个用户的ID
	IdLock.Lock()

	id := UidGen
	UidGen++

	IdLock.Unlock()

	//创建一个用户对象
	p := &User{
		Uid:  id,
		Conn: conn,
	}

	return p
}

func (u *User) GetId() int32 {
	return u.Uid
}

func (u *User) GetConn() iface.IConnection {
	return u.Conn
}

/*
* 提供一个发送给客户端消息的方法
* 主要是将pb的probuf数据序列化之后，再调用server的SendMsg
 */
func (u *User) SendMsg(msgId uint32, data proto.Message) {
	//将proto Message结构体序列化  转成二进制
	msg, err := proto.Marshal(data)
	if err != nil {
		fmt.Println("marshal msg err, err")
		return
	}

	//将二进制文件 通过server框架的sendMsg将数据发送给客户端
	if u.Conn == nil {
		fmt.Println("connection in user is nil")
		return
	}

	if err := u.Conn.SendMsg(msgId, msg); err != nil {
		fmt.Println("User SendMsg error!")
		return
	}

	return
}

// 同步用户的id
func (u *User) SyncUid() {
	//组件MsgID：0的proto数据
	proto_msg := &pb.SyncUid{
		Uid: u.Uid,
	}

	//发送消息给客户端
	u.SendMsg(0, proto_msg)
}

// 将用户拉取的持仓信息同步发送给客户端
func (u *User) SyncHoldingInfo(cid string) {
	//组建MsgId：200的proto数据

	//先查询cid客户的所有持仓记录
	rocords := GlobalHoldingRecord.PullRecordByCid(cid)

	//构建回送的proto消息
	var proto_record []*pb.AckHoldingInfo
	for _, tmp := range rocords {
		s := &pb.AckHoldingInfo{
			ClientId:     tmp.GetCliendId(),
			ExchangeType: tmp.GetExchangeType(),
			StockCode:    tmp.GetStockCode(),
			LastPrice:    tmp.GetLastPrice(),
			HoldAmount:   tmp.GetHoldAmount(),
			MarketValue:  tmp.GetMarketValue(),
		}

		proto_record = append(proto_record, s)
	}

	proto_msg := &pb.AckHoldingInfos{
		HInfo: proto_record,
	}

	//发送消息给客户端
	u.SendMsg(200, proto_msg)
}

// 执行用户的买入业务,
func (u *User) BuyingStock(cid string, et string, sc string, num int32) {
	//获取对应的股票
	var flag int32
	flag = 0
	s := GlobalStocks.GetStock(et, sc)
	if s != nil {
		//更新的操作
		flag = GlobalHoldingRecord.UpdataRecord(cid, *s, num, true)
	}

	BackToClient(u, flag, cid, *s, num)
}

// 执行用户的卖出业务,
func (u *User) SellingStock(cid string, et string, sc string, num int32) {
	//解析用户的请求
	// cid := proto_msg.GetCliendId()
	// et := proto_msg.GetExchangeType()
	// sc := proto_msg.GetStockCode()
	// num := proto_msg.GetNumber()

	//获取对应的股票
	var flag int32
	flag = 0
	s := GlobalStocks.GetStock(et, sc)
	if s != nil {
		//有对应的值
		//更新的操作
		flag = GlobalHoldingRecord.UpdataRecord(cid, *s, num, false)
	}

	BackToClient(u, flag, cid, *s, num)
}

// 构造执行结构回送给客户端，MsgID = 201
// 没有这只股票（0），新增记录（1），增加值（2），减少值（3），
// 删除记录（4），没有持股可以卖（5），没有足量的持股（6）
func BackToClient(u *User, flag int32, cid string, st Stock, num int32) {
	//组建MsgId：201的proto数据
	//构建回送的proto消息
	proto_msg := &pb.ChangeStock{}

	switch {
	case flag == 0:
		proto_msg = &pb.ChangeStock{
			RstType: pb.ResultType_NOSTOCK,
		}
	case flag == 1:
		proto_msg = &pb.ChangeStock{
			RstType:      pb.ResultType_ADDSTOCK,
			ExchangeType: st.GetExchangeType(),
			StockCode:    st.GetStockCode(),
			LastPrice:    st.GetLastPrice(),
			Number:       num,
		}
	case flag == 2:
		proto_msg = &pb.ChangeStock{
			RstType:      pb.ResultType_ADDVALUE,
			ExchangeType: st.GetExchangeType(),
			StockCode:    st.GetStockCode(),
			LastPrice:    st.GetLastPrice(),
			Number:       num,
		}
	case flag == 3:
		proto_msg = &pb.ChangeStock{
			RstType:      pb.ResultType_DECVALUE,
			ExchangeType: st.GetExchangeType(),
			StockCode:    st.GetStockCode(),
			LastPrice:    st.GetLastPrice(),
			Number:       num,
		}
	case flag == 4:
		proto_msg = &pb.ChangeStock{
			RstType:      pb.ResultType_DELSTOCK,
			ExchangeType: st.GetExchangeType(),
			StockCode:    st.GetStockCode(),
			LastPrice:    st.GetLastPrice(),
			Number:       num,
		}
	case flag == 5:
		proto_msg = &pb.ChangeStock{
			RstType:      pb.ResultType_NORECORD,
			ExchangeType: st.GetExchangeType(),
			StockCode:    st.GetStockCode(),
			LastPrice:    st.GetLastPrice(),
			Number:       num,
		}
	case flag == 6:
		proto_msg = &pb.ChangeStock{
			RstType:      pb.ResultType_NOEHOUGH,
			ExchangeType: st.GetExchangeType(),
			StockCode:    st.GetStockCode(),
			LastPrice:    st.GetLastPrice(),
			Number:       num,
		}
	}

	//发送消息给客户端
	u.SendMsg(201, proto_msg)
}

// 构造执行结构回送给客户端，MsgID = 202
func (u *User) SyncHistory(cid string) {
	data := ReadLog(cid)

	//组件MsgID：0的proto数据
	proto_msg := &pb.HistoryRecord{
		HisRcd: string(data),
	}

	//发送消息给客户端
	u.SendMsg(202, proto_msg)
}
