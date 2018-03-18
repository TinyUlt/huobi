package models

type KLineData struct {
	ID     int64   `json:"id"`     // K线ID
	Amount float64 `json:"amount"` // 成交量
	Count  int64   `json:"count"`  // 成交笔数
	Open   float64 `json:"open"`   // 开盘价
	Close  float64 `json:"close"`  // 收盘价, 当K线为最晚的一根时, 时最新成交价
	Low    float64 `json:"low"`    // 最低价
	High   float64 `json:"high"`   // 最高价
	Vol    float64 `json:"vol"`    // 成交额, 即SUM(每一笔成交价 * 该笔的成交数量)
}

type KLineReturn struct {
	Status  string      `json:"status"`   // 请求处理结果, "ok"、"error"
	Ts      int64       `json:"ts"`       // 响应生成时间点, 单位毫秒
	Data    []KLineData `json:"data"`     // KLine数据
	Ch      string      `json:"ch"`       // 数据所属的Channel, 格式: market.$symbol.kline.$period
	ErrCode string      `json:"err-code"` // 错误代码
	ErrMsg  string      `json:"err-msg"`  // 错误提示
}

type USDTData struct{
	Price float32 `json:"price"`
	MinTradeLimit float32 `json:"minTradeLimit"`
	MaxTradeLimit float32 `json:"maxTradeLimit"`
	PayMethod string `json:"payMethod"`
	UserName string `json:"userName"`
	TradeMonthTimes int32 `json:"tradeMonthTimes"`
}
type USDTReturn struct {
	TradeType string `json:"tradeType"`
	Data []USDTData `json:"data"`
}
type USDTPriceData struct {

	CoinId int `json:"coinId"`
	Price string `json:"price"`
}
type USDTPriceReturn struct{
	Data []USDTPriceData `json:"data"`
}

type ConfigData struct {
	CoinId string `json:"coinId"`
	TradeFollowSellEmial []string    `json:"tradeFollowSellEmail"`
	TradeFollowBuyEmial []string    `json:"tradeFollowBuyEmail"`
	TradeFollowStepEmail []string 		`json:"tradeFollowStepEmail"`
	Step float64    `json:"step"`
	PersonFollowName string    `json:"personFollowName"`
	PersonFollowSellEmail []string    `json:"personFollowSellEmail"`
	PersonFollowBuyEmail []string    `json:"personFollowBuyEmail"`
	AveFollowSellEmail []string    `json:"aveFollowSellEmail"`
	AveFollowBuyEmail []string    `json:"aveFollowBuyEmail"`
	PageCount int    `json:"pageCount"`
}