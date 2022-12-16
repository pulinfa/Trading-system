package iface

/*
* 封包和拆包的模块
* 直接面向TCP连接中的数据流，用于处理TCP粘包问题
 */
// import (
// 	"server/iface"
// )

type IDataPack interface {
	//获取包的header长度的方法
	GetHeadLen() uint32

	//封包的方法
	Pack(msg IMessage) ([]byte, error)

	//拆包的方法
	Unpack([]byte) (IMessage, error)
}
