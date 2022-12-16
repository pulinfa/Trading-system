package core

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

/*
* 这是日志的模块，为全局提供一个日志写的功能，和日志全局查某个用户的功能
* 不同的用户账号写不同的日志，并使用客户号进行命名
 */

// 交易结构体：时间，类型，市场，代码，数量
type Log struct {
	TimeStamp    time.Time
	Tp           int32 //主要分为，增加（1），减少（2），新建（3），删除（4）
	ClientId     string
	ExchangeType string
	StockCode    string
	Number       int32
}

func WriteLog(log Log) {
	filepath := "log/" + log.ClientId + ".txt"

	file, err := os.OpenFile(filepath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("open file err, ", err)
	}
	//及时关闭file句柄
	defer file.Close()

	s := "add"
	if log.Tp == 2 {
		s = "dec"
	} else if log.Tp == 3 {
		s = "new"
	} else {
		s = "delete"
	}

	//数据准备
	data := log.TimeStamp.Format("2006-01-02 15:04:05.000 Mon Jan") + " " + s + " " + log.ClientId + " " + log.ExchangeType + " " + log.StockCode + " " + fmt.Sprint(log.Number) + " \r\n"

	//写入文件时，使用带缓存的 *Writer
	write := bufio.NewWriter(file)
	write.WriteString(data)

	//Flush将缓存的文件真正写入到文件中
	write.Flush()
}

func ReadLog(cid string) []byte {
	filepath := "log/" + cid + ".txt"

	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Println(err)
	}

	// if it was successful in reading the file then
	// print out the contents as a string
	return data
}
