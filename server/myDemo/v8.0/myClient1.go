package main

import (
	"fmt"
	"io"
	"net"
	"time"

	"server/snet"
)

func main() {
	fmt.Println("client1 start")

	time.Sleep(1 * time.Second)

	//1 直接连接远程服务器，得到一个conn连接
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("client start err, exit")
		return
	}

	for {
		//发送封包的message信息
		dp := snet.NewDataPack()
		binaryMsg, err := dp.Pack(snet.NewMsgPackage(1, []byte("V8.0 client1 Test Message")))
		if err != nil {
			fmt.Println("Pack error: ", err)
			break
		}

		if _, err := conn.Write(binaryMsg); err != nil {
			fmt.Println("write error: ", err)
			break
		}

		//服务器应该给我们回复一个message数据， MsgID: 1 ping ping ping

		//先读取流中的head部分，得到ID和datalen
		binaryHead := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(conn, binaryHead); err != nil {
			fmt.Println("read head err: ", err)
			break
		}

		//将二进制的head拆包到msg结构体中
		msgHead, err := dp.Unpack(binaryHead)
		if err != nil {
			fmt.Println("client unpack msgHead error: ", err)
			break
		}

		//再根据datalen进行第二次读取，将data读出来
		if msgHead.GetMsgLen() > 0 {
			msg := msgHead.(*snet.Message)
			msg.Data = make([]byte, msg.GetMsgLen())

			if _, err := io.ReadFull(conn, msg.Data); err != nil {
				fmt.Println("read msg data error: ", err)
				return
			}

			fmt.Println("---> Recv Server Msg: ID = ", msg.Id, ", len = ", msg.DataLen, ", data = ", string(msg.Data))
		}

		time.Sleep(2 * time.Second)
	}
}
