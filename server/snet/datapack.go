package snet

/*
* 封包和拆包的具体模块
* 直接面向TCP连接中的数据流，用于处理TCP粘包问题
 */

import (
	"bytes"
	"encoding/binary"
	"errors"
	"server/iface"
	"server/utils"
)

type DataPack struct{}

// 拆包封包实例的初始化方法
func NewDataPack() *DataPack {
	return &DataPack{}
}

// 获取包的header长度的方法
func (dp *DataPack) GetHeadLen() uint32 {
	//DataLen uint32 (4字节) + ID uint32 (4字节)
	return 8
}

// 封包的方法
func (dp *DataPack) Pack(msg iface.IMessage) ([]byte, error) {
	//创建一个存放bytes字节的缓冲
	dataBuff := bytes.NewBuffer([]byte{})

	//将dataLen写进dataBuff中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgLen()); err != nil {
		return nil, err
	}
	//将ID写进dataBuff中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}

	//将data数据写入dataBuff中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}

	return dataBuff.Bytes(), nil
}

// 拆包的方法
// 1. 先将header的信息读出来，
// 2. 根据header中的信息里面的data的长度，再进行一次读
func (dp *DataPack) Unpack(binaryData []byte) (iface.IMessage, error) {
	//创建一个从输入的二进制数据的ioReader
	dataBuff := bytes.NewReader(binaryData)

	//只解压header信息，得到datalen和id
	msg := &Message{}

	//读取datalen
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}
	//读取MsgID
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}

	//判断datalen是否已经超出了我们允许的最大包长度
	if utils.GlobalObject.MaxPackageSize > 0 && msg.DataLen > utils.GlobalObject.MaxPackageSize {
		return nil, errors.New("too Large msg data recv")
	}

	return msg, nil
}
