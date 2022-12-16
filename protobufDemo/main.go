package main

import (
	"fmt"
	"server/snet"

	"demo/pb"

	"github.com/golang/protobuf/proto"
)

func main() {
	//定义一个Person结构对象
	person := &pb.Person{
		Name:   "pu linfa",
		Id:     2020229019,
		Emails: []string{"abc@qq.com, bcd@tju.edu.com, cde@gmail.com"},
		Phones: []*pb.PhoneNumber{
			&pb.PhoneNumber{
				Number: "1111",
				Type:   pb.PhoneType_HOME,
			},
			&pb.PhoneNumber{
				Number: "2222",
				Type:   pb.PhoneType_MOBILE,
			},
			&pb.PhoneNumber{
				Number: "3333",
				Type:   pb.PhoneType_WORK,
			},
		},
	}

	//将person对象  就是将protobuf的message进行序列化，得到一个二进制文件
	data, err := proto.Marshal(person)
	if err != nil {
		fmt.Println("marshal err, ", err)
	}

	//解码
	newdata := &pb.Person{}
	err = proto.Unmarshal(data, newdata)
	if err != nil {
		fmt.Println("unmarshal err, ", err)
	}

	fmt.Println("源数据：", person)
	fmt.Println("序列化数据： ", data)
	fmt.Println("解码之后的数据：", newdata)

	s := snet.NewServer("version:10.0")

	s.Start()
}
