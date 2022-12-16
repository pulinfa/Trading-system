package core

/*
* 股票模块：对于每个股票来讲，（市场，股票）----> (价格)
 */
type Stock struct {
	exchangeType string  //市场，1表示沪，2表示深
	stockCode    string  //代码	沪市：600000-600099，深市：000001-000100
	lastPrice    float64 //最新价格 10.00-1000.00之间的随机数
}

// 新建一只股票
func NewStock(ET string, SC string, LP float64) *Stock {
	s := &Stock{
		exchangeType: ET,
		stockCode:    SC,
		lastPrice:    LP,
	}

	return s
}

// 取股票的市场类型
func (st *Stock) GetExchangeType() string {
	return st.exchangeType
}

// 取股票的代码
func (st *Stock) GetStockCode() string {
	return st.stockCode
}

// 取股票的最新价格
func (st *Stock) GetLastPrice() float64 {
	return st.lastPrice
}

// 设置股票的最新价格
func (st *Stock) SetLastPrice(lp float64) {
	st.lastPrice = lp
}

// 所有股票信息的集合
type Stocks struct {
	stocks []*Stock
}

// 创建一个新的股票仓库
func NewStocks() *Stocks {
	s := &Stocks{
		stocks: make([]*Stock, 0),
	}

	return s
}

func (sts *Stocks) GetStocks() []*Stock {
	QueryAllStock()
	return sts.stocks
}

// 在股票仓库中新增加一只股票
func (sts *Stocks) AddStock(st *Stock) {
	sts.stocks = append(sts.stocks, st)
}

// 在股票仓库中查找一只股票，使用市场和股票代码进行查询
func (sts *Stocks) GetStock(ET string, SC string) *Stock {
	//使用数据库进行查询

	for i := 0; i < len(sts.stocks); i++ {
		if sts.stocks[i].GetExchangeType() == ET && sts.stocks[i].GetStockCode() == SC {
			return sts.stocks[i]
		}
	}

	return nil
}

// 定义一个全局的股票仓库
var GlobalStocks *Stocks

// // 提供一个初始化方法，初始化当前的GlobalStocks
func InitStocks() {
	GlobalStocks = NewStocks()

	QueryAllStock()

	/* //添加几个新的股票
	s1 := NewStock("1", "600001", 58.32)
	s2 := NewStock("1", "600002", 121.68)
	s3 := NewStock("2", "000001", 69.45)
	s4 := NewStock("2", "000002", 131.14)

	GlobalStocks.AddStock(s1)
	GlobalStocks.AddStock(s2)
	GlobalStocks.AddStock(s3)
	GlobalStocks.AddStock(s4) */
}
