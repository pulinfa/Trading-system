package apis

import (
	"fmt"
	"server/iface"
	"server/snet"
	"trasys/core"
	"trasys/pb"

	"github.com/golang/protobuf/proto"
)

/*
* 这是卖出股票的方法实现
 */

type Selling struct {
	snet.BaseRouter
}

func (selling *Selling) Handle(request iface.IRequest) {
	//1. 解析客户端传递过来的请求
	proto_msg := &pb.TradingInfo{}
	err := proto.Unmarshal(request.GetData(), proto_msg)
	if err != nil {
		fmt.Println("pull Unmarshal error, ", err)
		return
	}

	//2. 当前的数据是由哪个客户端发送的
	uid, err := request.GetConnection().GetProperty("uid")
	if err != nil {
		fmt.Println("no such property, ", err)
	}

	//3. 得到对应的user对象
	user := core.GlobalUserMgr.GetUserByUid(uid.(int32))

	//5. 发送回客户端

	cid := proto_msg.GetClientId()
	et := proto_msg.GetExchangeType()
	sc := proto_msg.GetStockCode()
	num := proto_msg.GetNumber()

	user.SellingStock(cid, et, sc, num)
}
