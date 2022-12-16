package core

import (
	"fmt"
	"time"
)

//持仓信息

// 持仓消息结构体
type HoldingInfo struct {
	clientId    string  //客户号 范围：000000000001-999999999999,总共12位
	stock       Stock   //所持有的股票
	holdAmount  int32   //持仓数量
	marketValue float64 //市值 holdAmount * stock.GetLastPrice()
}

// 创建一个新的持仓，这个的话主要的作用是说，客户之前没有这只股票，所以这个时候需要新创建一个持股的情况
func NewHoldingInfo(clientid string, st Stock, holdamount int32) *HoldingInfo {
	s := &HoldingInfo{
		clientId:    clientid,
		stock:       st,
		holdAmount:  holdamount,
		marketValue: float64(holdamount) * st.GetLastPrice(),
	}

	return s
}

// 得到持仓的客户号
func (hi *HoldingInfo) GetCliendId() string {
	return hi.clientId
}

// 得到持仓的股票
func (hi *HoldingInfo) GetStock() Stock {
	return hi.stock
}

// 得到持仓的数量
func (hi *HoldingInfo) GetHoldAmount() int32 {
	return hi.holdAmount
}

// 得到持仓的市值
func (hi *HoldingInfo) GetMarketValue() float64 {
	return hi.marketValue
}

// 取股票的市场类型
func (hi *HoldingInfo) GetExchangeType() string {
	return hi.stock.GetExchangeType()
}

// 取股票的代码
func (hi *HoldingInfo) GetStockCode() string {
	return hi.stock.GetStockCode()
}

// 取股票的最新价格
func (hi *HoldingInfo) GetLastPrice() float64 {
	return hi.stock.GetLastPrice()
}

// 增加持仓数量，就是所谓的加仓
func (hi *HoldingInfo) IncHoldAmount(num int32) {
	hi.holdAmount += num
	hi.marketValue = hi.stock.GetLastPrice() * float64(hi.holdAmount)
}

// 减少持仓数量，就是所谓的卖出
func (hi *HoldingInfo) DecHoldAmount(num int32) int32 {
	//这个时候卖出的数量不够，没有那么多股票可以卖出
	if num > hi.holdAmount {
		fmt.Println("sold num is too much, you don't have enough stock to sell,")
		return 6
	}
	hi.holdAmount -= num
	hi.marketValue = hi.stock.GetLastPrice() * float64(hi.holdAmount)

	return 3
}

// 持仓情况记录表，记录了所有的持仓情况
type HoldingRecord struct {
	holdingRecord []*HoldingInfo
}

// 持仓记录表的初始化
func NewHoldingRecord() *HoldingRecord {
	s := &HoldingRecord{
		holdingRecord: make([]*HoldingInfo, 0),
	}

	return s
}

func (hr *HoldingRecord) GetHodingRecord() []*HoldingInfo {
	return hr.holdingRecord
}

// 增加一条记录
func (hr *HoldingRecord) AddRecord(hi *HoldingInfo) {
	hr.holdingRecord = append(hr.holdingRecord, hi)
}

// 删除一条记录
func (hr *HoldingRecord) DelRecord(hi *HoldingInfo) {
	j := 0
	for _, v := range hr.holdingRecord {
		if v != hi {
			hr.holdingRecord[j] = v
			j++
		}
	}

	hr.holdingRecord = hr.holdingRecord[:j]
}

// 查找持仓记录，使用客户号和股票信息
func (hr *HoldingRecord) SearchRecord(ci string, st Stock) *HoldingInfo {
	for _, v := range hr.holdingRecord {
		if v.clientId == ci && v.stock.GetExchangeType() == st.GetExchangeType() && v.stock.GetStockCode() == st.GetStockCode() {
			return v
		}
	}

	return nil
}

// 通过用户的Cid查询所有相关的持仓记录
func (hr *HoldingRecord) PullRecordByCid(cid string) []*HoldingInfo {
	var ans []*HoldingInfo

	for _, tmp := range hr.holdingRecord {
		if tmp.clientId == cid {
			ans = append(ans, tmp)
		}
	}

	return ans
}

// 更新持仓记录，分别包含了几种情况，根据类型号的不同进行不同的操作，使用客户号和股票信息
// 1. 增加的操作，类型为true
// 2. 减少的操作，类型为false
// 返回值：新增记录（1），增加值（2），减少值（3），删除记录（4），没有持股可以卖（5），没有足量的持股（6）
func (hr *HoldingRecord) UpdataRecord(ci string, st Stock, number int32, tp bool) int32 {
	if tp {
		//增加持仓
		tmp := hr.SearchRecord(ci, st)
		if tmp == nil {
			fmt.Println("this Record is not FOUND")
			hr.AddRecord(NewHoldingInfo(ci, st, number))
			//写日志
			s := Log{
				TimeStamp:    time.Now(),
				Tp:           3,
				ClientId:     ci,
				ExchangeType: st.GetExchangeType(),
				StockCode:    st.GetStockCode(),
				Number:       number,
			}
			WriteLog(s)

			//想数据库中Insert一条新的记录
			InsertRecord(ci, st.GetExchangeType(), st.GetStockCode(), st.GetLastPrice(), number, float64(number)*st.GetLastPrice())

			return 1
		} else {
			fmt.Println("FOUND it")
			tmp.IncHoldAmount(number)
			//写日志
			s := Log{
				TimeStamp:    time.Now(),
				Tp:           1,
				ClientId:     ci,
				ExchangeType: st.GetExchangeType(),
				StockCode:    st.GetStockCode(),
				Number:       number,
			}
			WriteLog(s)

			//在数据库中Updata一条数据
			UpdateRecord(ci, tmp.GetExchangeType(), tmp.GetStockCode(), tmp.GetHoldAmount()+number, float64(tmp.GetHoldAmount()+number)*st.GetLastPrice())

			return 2
		}

	} else {
		//减少持仓
		tmp := hr.SearchRecord(ci, st)
		if tmp == nil {
			fmt.Println("this Record is not FOUND, so Dec a nil record is not ledge")
			return 5
		} else {
			fmt.Println("FOUND it")
			if tmp.GetHoldAmount() != number {
				ans := tmp.DecHoldAmount(number)

				if ans == 3 {
					//写日志
					s := Log{
						TimeStamp:    time.Now(),
						Tp:           2,
						ClientId:     ci,
						ExchangeType: st.GetExchangeType(),
						StockCode:    st.GetStockCode(),
						Number:       number,
					}
					WriteLog(s)

					//数据库更新
					UpdateRecord(ci, tmp.GetExchangeType(), tmp.GetStockCode(), tmp.GetHoldAmount()-number, float64(tmp.GetHoldAmount()-number)*st.GetLastPrice())
				}

				return ans
			} else {
				hr.DelRecord(tmp)
				//写日志
				s := Log{
					TimeStamp:    time.Now(),
					Tp:           4,
					ClientId:     ci,
					ExchangeType: st.GetExchangeType(),
					StockCode:    st.GetStockCode(),
					Number:       number,
				}
				WriteLog(s)

				//delete
				DeleteRecord(ci, tmp.GetExchangeType(), tmp.GetStockCode())

				return 4
			}
		}
	}

	return 0
}

// // 定义一个全局的持仓记录表
var GlobalHoldingRecord *HoldingRecord

// // 提供一个初始化方法，初始化当前的GlobalHoldingRecord
func init() {
	//初始化一个股票仓库
	InitDB()
	InitStocks()

	GlobalHoldingRecord = NewHoldingRecord()

	QueryAllReocrd()

	/* //需要添加一些持仓记录
	s1 := NewHoldingInfo("000000000001", *GlobalStocks.stocks[0], 1000)
	s2 := NewHoldingInfo("000000000002", *GlobalStocks.stocks[1], 999)
	s3 := NewHoldingInfo("000000000003", *GlobalStocks.stocks[2], 1670)
	s4 := NewHoldingInfo("000000000004", *GlobalStocks.stocks[3], 190)

	GlobalHoldingRecord.AddRecord(s1)
	GlobalHoldingRecord.AddRecord(s2)
	GlobalHoldingRecord.AddRecord(s3)
	GlobalHoldingRecord.AddRecord(s4) */
}
