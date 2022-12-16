package apis

/*
* 历史记录查询
 */

import (
	"fmt"
	"server/iface"
	"server/snet"

	"trasys/core"
	"trasys/pb"

	"github.com/golang/protobuf/proto"
)

type PullHistoryRecords struct {
	snet.BaseRouter
}

func (history *PullHistoryRecords) Handle(request iface.IRequest) {
	//1. 解析客户端传递过来的请求
	proto_msg := &pb.PullHistory{}
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

	//4. 从信息中提取出客户号
	cid := proto_msg.GetCid()

	//5. 发送回客户端
	user.SyncHistory(cid)
}
