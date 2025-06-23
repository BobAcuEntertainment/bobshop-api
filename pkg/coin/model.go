package coin

import (
	"app/pkg/ccxt"
	"app/pkg/enum"
	"math"

	"github.com/nhnghia272/gopkg"
)

type SymbolQuery struct {
	Sandbox bool `json:"sandbox" form:"sandbox" validate:"omitempty"`
}

type TimeframeQuery struct {
	Sandbox bool `json:"sandbox" form:"sandbox" validate:"omitempty"`
}

type CurrencyQuery struct {
	Sandbox bool `json:"sandbox" form:"sandbox" validate:"omitempty"`
}

type Currency struct {
	Id        string             `json:"currency_id"`
	Code      string             `json:"code"`
	Precision float64            `json:"precision"`
	Name      string             `json:"name"`
	Fee       float64            `json:"fee"`
	Active    bool               `json:"active"`
	Deposit   bool               `json:"deposit"`
	Withdraw  bool               `json:"withdraw"`
	NumericId int64              `json:"numeric_id"`
	Type      string             `json:"type"`
	Margin    bool               `json:"margin"`
	Limits    CurrencyLimits     `json:"currency_limits"`
	Networks  map[string]Network `json:"networks"`
}

func NewCurrency(data ccxt.Currency) *Currency {
	currency := &Currency{
		Id:        gopkg.Value(data.Id),
		Code:      gopkg.Value(data.Code),
		Precision: gopkg.Value(data.Precision),
		Name:      gopkg.Value(data.Name),
		Fee:       gopkg.Value(data.Fee),
		Active:    gopkg.Value(data.Active),
		Deposit:   gopkg.Value(data.Deposit),
		Withdraw:  gopkg.Value(data.Withdraw),
		NumericId: gopkg.Value(data.NumericId),
		Type:      gopkg.Value(data.Type),
		Margin:    gopkg.Value(data.Margin),
		Limits:    NewCurrencyLimits(data.Limits),
	}
	if math.IsNaN(currency.Precision) {
		currency.Precision = 0
	}
	if math.IsNaN(currency.Fee) {
		currency.Fee = 0
	}
	currency.Networks = make(map[string]Network)
	for _, network := range data.Networks {
		currency.Networks[gopkg.Value(network.Id)] = NewNetwork(network)
	}

	return currency
}

type MarketQuery struct {
	Sandbox    bool            `json:"sandbox" form:"sandbox" validate:"omitempty"`
	MarketType enum.MarketType `json:"market_type" form:"market_type" validate:"required,market_type"`
}

type Market struct {
	ID     string `json:"market_id"`
	Symbol string `json:"symbol"`
	// UppercaseId    string    `json:"uppercase_id"`
	// LowercaseId    string    `json:"lowercase_id"`
	// BaseCurrency   string    `json:"base_currency"`
	// QuoteCurrency  string    `json:"quote_currency"`
	// BaseId         string    `json:"base_id"`
	// QuoteId        string    `json:"quote_id"`
	// Active         bool      `json:"active"`
	Type   string `json:"type"`
	Spot   bool   `json:"spot"`
	Margin bool   `json:"margin"`
	// Swap           bool      `json:"swap"`
	// Future         bool      `json:"future"`
	// Option         bool      `json:"option"`
	// Contract       bool      `json:"contract"`
	// Settle         string    `json:"settle"`
	// SettleId       string    `json:"settle_id"`
	// ContractSize   float64   `json:"contract_size"`
	// Linear         bool      `json:"linear"`
	// Inverse        bool      `json:"inverse"`
	// Quanto         bool      `json:"quanto"`
	// Expiry         int64     `json:"expiry"`
	// ExpiryDatetime string    `json:"expiry_datetime"`
	// Strike         float64   `json:"strike"`
	// OptionType     string    `json:"option_type"`
	// Taker          float64   `json:"taker"`
	// Maker          float64   `json:"maker"`
	// Limits         Limits    `json:"limits"`
	// Precision      Precision `json:"precision"`
	// Created        int64     `json:"created"`
}

func NewMarket(data ccxt.MarketInterface) *Market {
	market := &Market{
		ID:     gopkg.Value(data.Symbol),
		Symbol: gopkg.Value(data.Symbol),
		// UppercaseId:    gopkg.Value(data.UppercaseId),
		// LowercaseId:    gopkg.Value(data.LowercaseId),
		// BaseCurrency:   gopkg.Value(data.BaseCurrency),
		// QuoteCurrency:  gopkg.Value(data.QuoteCurrency),
		// BaseId:         gopkg.Value(data.BaseId),
		// QuoteId:        gopkg.Value(data.QuoteId),
		// Active:         gopkg.Value(data.Active),
		Type:   gopkg.Value(data.Type),
		Spot:   gopkg.Value(data.Spot),
		Margin: gopkg.Value(data.Margin),
		// Swap:           gopkg.Value(data.Swap),
		// Future:         gopkg.Value(data.Future),
		// Option:         gopkg.Value(data.Option),
		// Contract:       gopkg.Value(data.Contract),
		// Settle:         gopkg.Value(data.Settle),
		// SettleId:       gopkg.Value(data.SettleId),
		// ContractSize:   gopkg.Value(data.ContractSize),
		// Linear:         gopkg.Value(data.Linear),
		// Inverse:        gopkg.Value(data.Inverse),
		// Quanto:         gopkg.Value(data.Quanto),
		// Expiry:         gopkg.Value(data.Expiry),
		// ExpiryDatetime: gopkg.Value(data.ExpiryDatetime),
		// Strike:         gopkg.Value(data.Strike),
		// OptionType:     gopkg.Value(data.OptionType),
		// Taker:          gopkg.Value(data.Taker),
		// Maker:          gopkg.Value(data.Maker),
		// Limits:         NewLimits(data.Limits),
		// Created:        gopkg.Value(data.Created),
	}
	// if math.IsNaN(market.ContractSize) {
	// 	market.ContractSize = 0
	// }
	// if math.IsNaN(market.Strike) {
	// 	market.Strike = 0
	// }
	// if math.IsNaN(market.Taker) {
	// 	market.Taker = 0
	// }
	// if math.IsNaN(market.Maker) {
	// 	market.Maker = 0
	// }
	// if precision, ok := data.Info["precision"]; ok {
	// 	market.Precision = NewPrecision(ccxt.NewPrecision(precision))
	// }
	return market
}

type TickerQuery struct {
	Sandbox bool     `json:"sandbox" form:"sandbox" validate:"omitempty"`
	Symbols []string `json:"symbols" form:"symbols" validate:"required,min=1"`
}

type Ticker struct {
	Symbol        string  `json:"symbol"`
	Timestamp     int64   `json:"timestamp"`
	Datetime      string  `json:"datetime"`
	High          float64 `json:"high"`
	Low           float64 `json:"low"`
	Bid           float64 `json:"bid"`
	BidVolume     float64 `json:"bid_volume"`
	Ask           float64 `json:"ask"`
	AskVolume     float64 `json:"ask_volume"`
	Vwap          float64 `json:"vwap"`
	Open          float64 `json:"open"`
	Close         float64 `json:"close"`
	Last          float64 `json:"last"`
	PreviousClose float64 `json:"previous_close"`
	Change        float64 `json:"change"`
	Percentage    float64 `json:"percentage"`
	Average       float64 `json:"average"`
	BaseVolume    float64 `json:"base_volume"`
	QuoteVolume   float64 `json:"quote_volume"`
}

func NewTicker(data ccxt.Ticker) *Ticker {
	ticker := &Ticker{
		Symbol:        gopkg.Value(data.Symbol),
		Timestamp:     gopkg.Value(data.Timestamp),
		Datetime:      gopkg.Value(data.Datetime),
		High:          gopkg.Value(data.High),
		Low:           gopkg.Value(data.Low),
		Bid:           gopkg.Value(data.Bid),
		BidVolume:     gopkg.Value(data.BidVolume),
		Ask:           gopkg.Value(data.Ask),
		AskVolume:     gopkg.Value(data.AskVolume),
		Vwap:          gopkg.Value(data.Vwap),
		Open:          gopkg.Value(data.Open),
		Close:         gopkg.Value(data.Close),
		Last:          gopkg.Value(data.Last),
		PreviousClose: gopkg.Value(data.PreviousClose),
		Change:        gopkg.Value(data.Change),
		Percentage:    gopkg.Value(data.Percentage),
		Average:       gopkg.Value(data.Average),
		BaseVolume:    gopkg.Value(data.BaseVolume),
		QuoteVolume:   gopkg.Value(data.QuoteVolume),
	}
	if math.IsNaN(ticker.High) {
		ticker.High = 0
	}
	if math.IsNaN(ticker.Low) {
		ticker.Low = 0
	}
	if math.IsNaN(ticker.Bid) {
		ticker.Bid = 0
	}
	if math.IsNaN(ticker.BidVolume) {
		ticker.BidVolume = 0
	}
	if math.IsNaN(ticker.Ask) {
		ticker.Ask = 0
	}
	if math.IsNaN(ticker.AskVolume) {
		ticker.AskVolume = 0
	}
	if math.IsNaN(ticker.Vwap) {
		ticker.Vwap = 0
	}
	if math.IsNaN(ticker.Open) {
		ticker.Open = 0
	}
	if math.IsNaN(ticker.Close) {
		ticker.Close = 0
	}
	if math.IsNaN(ticker.Last) {
		ticker.Last = 0
	}
	if math.IsNaN(ticker.PreviousClose) {
		ticker.PreviousClose = 0
	}
	if math.IsNaN(ticker.Change) {
		ticker.Change = 0
	}
	if math.IsNaN(ticker.Percentage) {
		ticker.Percentage = 0
	}
	if math.IsNaN(ticker.Average) {
		ticker.Average = 0
	}
	if math.IsNaN(ticker.BaseVolume) {
		ticker.BaseVolume = 0
	}
	if math.IsNaN(ticker.QuoteVolume) {
		ticker.QuoteVolume = 0
	}
	return ticker
}

type Precision struct {
	Amount float64 `json:"amount"`
	Price  float64 `json:"price"`
}

func NewPrecision(data ccxt.Precision) Precision {
	precision := Precision{
		Amount: gopkg.Value(data.Amount),
		Price:  gopkg.Value(data.Price),
	}
	if math.IsNaN(precision.Amount) {
		precision.Amount = 0
	}
	if math.IsNaN(precision.Price) {
		precision.Price = 0
	}
	return precision
}

type MinMax struct {
	Min float64 `json:"min"`
	Max float64 `json:"max"`
}

func NewMinMax(data ccxt.MinMax) MinMax {
	minmax := MinMax{
		Min: gopkg.Value(data.Min),
		Max: gopkg.Value(data.Max),
	}
	if math.IsNaN(minmax.Min) {
		minmax.Min = 0
	}
	if math.IsNaN(minmax.Max) {
		minmax.Max = 0
	}
	return minmax
}

type CurrencyLimits struct {
	Amount   MinMax `json:"amount"`
	Withdraw MinMax `json:"withdraw"`
}

func NewCurrencyLimits(data ccxt.CurrencyLimits) CurrencyLimits {
	return CurrencyLimits{
		Amount:   NewMinMax(data.Amount),
		Withdraw: NewMinMax(data.Withdraw),
	}
}

type Limits struct {
	Amount   MinMax `json:"amount"`
	Cost     MinMax `json:"cost"`
	Leverage MinMax `json:"leverage"`
	Price    MinMax `json:"price"`
}

func NewLimits(data ccxt.Limits) Limits {
	return Limits{
		Amount:   NewMinMax(data.Amount),
		Cost:     NewMinMax(data.Cost),
		Leverage: NewMinMax(data.Leverage),
		Price:    NewMinMax(data.Price),
	}
}

type Network struct {
	Id        string         `json:"network_id"`
	Fee       float64        `json:"fee"`
	Active    bool           `json:"active"`
	Deposit   bool           `json:"deposit"`
	Withdraw  bool           `json:"withdraw"`
	Precision float64        `json:"precision"`
	Limits    CurrencyLimits `json:"currency_limits"`
}

func NewNetwork(data ccxt.Network) Network {
	network := Network{
		Id:        gopkg.Value(data.Id),
		Fee:       gopkg.Value(data.Fee),
		Active:    gopkg.Value(data.Active),
		Deposit:   gopkg.Value(data.Deposit),
		Withdraw:  gopkg.Value(data.Withdraw),
		Precision: gopkg.Value(data.Precision),
		Limits:    NewCurrencyLimits(data.Limits),
	}
	if math.IsNaN(network.Fee) {
		network.Fee = 0
	}
	return network
}

type OHLCVQuery struct {
	Sandbox   bool   `json:"sandbox" form:"sandbox" validate:"omitempty"`
	Symbol    string `json:"symbol" form:"symbol" validate:"required"`
	Timeframe string `json:"timeframe" form:"timeframe" validate:"required"`
}

type OHLCV struct {
	Timestamp int64   `json:"timestamp"`
	Open      float64 `json:"open"`
	High      float64 `json:"high"`
	Low       float64 `json:"low"`
	Close     float64 `json:"close"`
	Volume    float64 `json:"volume"`
}

func NewOHLCV(data ccxt.OHLCV) *OHLCV {
	ohlcv := &OHLCV{
		Timestamp: data.Timestamp,
		Open:      data.Open,
		High:      data.High,
		Low:       data.Low,
		Close:     data.Close,
		Volume:    data.Volume,
	}
	if math.IsNaN(ohlcv.Open) {
		ohlcv.Open = 0
	}
	if math.IsNaN(ohlcv.High) {
		ohlcv.High = 0
	}
	if math.IsNaN(ohlcv.Low) {
		ohlcv.Low = 0
	}
	if math.IsNaN(ohlcv.Close) {
		ohlcv.Close = 0
	}
	if math.IsNaN(ohlcv.Volume) {
		ohlcv.Volume = 0
	}
	return ohlcv
}

type TradeQuery struct {
	Sandbox bool   `json:"sandbox" form:"sandbox" validate:"omitempty"`
	Symbol  string `json:"symbol" form:"symbol" validate:"required"`
}

type Trade struct {
	Amount       float64 `json:"amount"`
	Price        float64 `json:"price"`
	Cost         float64 `json:"cost"`
	Id           string  `json:"id"`
	Order        string  `json:"order"`
	Timestamp    int64   `json:"timestamp"`
	Datetime     string  `json:"datetime"`
	Symbol       string  `json:"symbol"`
	Type         string  `json:"type"`
	Side         string  `json:"side"`
	TakerOrMaker string  `json:"taker_or_maker"`
	Fee          Fee     `json:"fee"`
}

func NewTrade(data ccxt.Trade) *Trade {
	trade := &Trade{
		Amount:       gopkg.Value(data.Amount),
		Price:        gopkg.Value(data.Price),
		Cost:         gopkg.Value(data.Cost),
		Id:           gopkg.Value(data.Id),
		Order:        gopkg.Value(data.Order),
		Timestamp:    gopkg.Value(data.Timestamp),
		Datetime:     gopkg.Value(data.Datetime),
		Symbol:       gopkg.Value(data.Symbol),
		Type:         gopkg.Value(data.Type),
		Side:         gopkg.Value(data.Side),
		TakerOrMaker: gopkg.Value(data.TakerOrMaker),
		Fee:          NewFee(data.Fee),
	}
	if math.IsNaN(trade.Amount) {
		trade.Amount = 0
	}
	if math.IsNaN(trade.Price) {
		trade.Price = 0
	}
	if math.IsNaN(trade.Cost) {
		trade.Cost = 0
	}
	return trade
}

type Fee struct {
	Rate float64 `json:"rate"`
	Cost float64 `json:"cost"`
}

func NewFee(data ccxt.Fee) Fee {
	fee := Fee{
		Rate: gopkg.Value(data.Rate),
		Cost: gopkg.Value(data.Cost),
	}
	if math.IsNaN(fee.Rate) {
		fee.Rate = 0
	}
	if math.IsNaN(fee.Cost) {
		fee.Cost = 0
	}
	return fee
}

type OrderBookQuery struct {
	Sandbox bool   `json:"sandbox" form:"sandbox" validate:"omitempty"`
	Symbol  string `json:"symbol" form:"symbol" validate:"required"`
}

type OrderBook struct {
	Bids      [][]float64 `json:"bids"`
	Asks      [][]float64 `json:"asks"`
	Symbol    string      `json:"symbol"`
	Timestamp int64       `json:"timestamp"`
	Datetime  string      `json:"datetime"`
	Nonce     int64       `json:"nonce"`
}

func NewOrderBook(data ccxt.OrderBook) *OrderBook {
	orderbook := &OrderBook{
		Bids:      data.Bids,
		Asks:      data.Asks,
		Symbol:    gopkg.Value(data.Symbol),
		Timestamp: gopkg.Value(data.Timestamp),
		Datetime:  gopkg.Value(data.Datetime),
		Nonce:     gopkg.Value(data.Nonce),
	}
	return orderbook
}

type BalanceQuery struct {
	Sandbox    bool            `json:"sandbox" form:"sandbox" validate:"omitempty"`
	TraderId   enum.Trader     `json:"trader_id" form:"trader_id" validate:"required,len=24"`
	MarketType enum.MarketType `json:"market_type" form:"market_type" validate:"required,market_type"`
}

type Balance struct {
	Code  string  `json:"code"`
	Free  float64 `json:"free"`
	Used  float64 `json:"used"`
	Total float64 `json:"total"`
}

func NewBalance(code string, data ccxt.Balance) *Balance {
	balance := &Balance{
		Code:  code,
		Free:  gopkg.Value(data.Free),
		Used:  gopkg.Value(data.Used),
		Total: gopkg.Value(data.Total),
	}
	if math.IsNaN(balance.Free) {
		balance.Free = 0
	}
	if math.IsNaN(balance.Used) {
		balance.Used = 0
	}
	if math.IsNaN(balance.Total) {
		balance.Total = 0
	}
	return balance
}

type OrderQuery struct {
	Sandbox  bool        `json:"sandbox" form:"sandbox" validate:"omitempty"`
	Symbol   string      `json:"symbol" form:"symbol" validate:"required"`
	TraderId enum.Trader `json:"trader_id" form:"trader_id" validate:"required,len=24"`
}

type Order struct {
	Id                 string   `json:"order_id"`
	ClientOrderId      string   `json:"client_order_id"`
	Timestamp          int64    `json:"timestamp"`
	Datetime           string   `json:"datetime"`
	LastTradeTimestamp string   `json:"last_trade_timestamp"`
	Symbol             string   `json:"symbol"`
	Type               string   `json:"type"`
	Side               string   `json:"side"`
	Price              float64  `json:"price"`
	Cost               float64  `json:"cost"`
	Average            float64  `json:"average"`
	Amount             float64  `json:"amount"`
	Filled             float64  `json:"filled"`
	Remaining          float64  `json:"remaining"`
	Status             string   `json:"status"`
	ReduceOnly         bool     `json:"reduce_only"`
	PostOnly           bool     `json:"post_only"`
	Fee                Fee      `json:"fee"`
	Trades             []*Trade `json:"trades"`
	TriggerPrice       float64  `json:"trigger_price"`
	StopLossPrice      float64  `json:"stop_loss_price"`
	TakeProfitPrice    float64  `json:"take_profit_price"`
}

func NewOrder(data ccxt.Order) *Order {
	order := &Order{
		Id:                 gopkg.Value(data.Id),
		ClientOrderId:      gopkg.Value(data.ClientOrderId),
		Timestamp:          gopkg.Value(data.Timestamp),
		Datetime:           gopkg.Value(data.Datetime),
		LastTradeTimestamp: gopkg.Value(data.LastTradeTimestamp),
		Symbol:             gopkg.Value(data.Symbol),
		Type:               gopkg.Value(data.Type),
		Side:               gopkg.Value(data.Side),
		Price:              gopkg.Value(data.Price),
		Cost:               gopkg.Value(data.Cost),
		Average:            gopkg.Value(data.Average),
		Amount:             gopkg.Value(data.Amount),
		Filled:             gopkg.Value(data.Filled),
		Remaining:          gopkg.Value(data.Remaining),
		Status:             gopkg.Value(data.Status),
		ReduceOnly:         gopkg.Value(data.ReduceOnly),
		PostOnly:           gopkg.Value(data.PostOnly),
		Fee:                NewFee(data.Fee),
		Trades:             make([]*Trade, 0),
		TriggerPrice:       gopkg.Value(data.TriggerPrice),
		StopLossPrice:      gopkg.Value(data.StopLossPrice),
		TakeProfitPrice:    gopkg.Value(data.TakeProfitPrice),
	}
	if math.IsNaN(order.Price) {
		order.Price = 0
	}
	if math.IsNaN(order.Cost) {
		order.Cost = 0
	}
	if math.IsNaN(order.Average) {
		order.Average = 0
	}
	if math.IsNaN(order.Amount) {
		order.Amount = 0
	}
	if math.IsNaN(order.Filled) {
		order.Filled = 0
	}
	if math.IsNaN(order.Remaining) {
		order.Remaining = 0
	}
	if math.IsNaN(order.TriggerPrice) {
		order.TriggerPrice = 0
	}
	if math.IsNaN(order.StopLossPrice) {
		order.StopLossPrice = 0
	}
	if math.IsNaN(order.TakeProfitPrice) {
		order.TakeProfitPrice = 0
	}
	order.Trades = gopkg.MapFunc(data.Trades, func(trade ccxt.Trade) *Trade { return NewTrade(trade) })
	return order
}

type OrderBody struct {
	Symbol    string         `json:"symbol"`
	OrderType enum.OrderType `json:"order_type"`
	OrderSide enum.OrderSide `json:"order_side"`
	Amount    float64        `json:"amount"`
	Price     float64        `json:"price"`
}

type OrderData struct {
	Sandbox   bool           `json:"sandbox" validate:"omitempty"`
	Symbol    string         `json:"symbol" validate:"required"`
	OrderType enum.OrderType `json:"order_type" validate:"required,order_type"`
	OrderSide enum.OrderSide `json:"order_side" validate:"required,order_side"`
	Amount    float64        `json:"amount" validate:"required"`
	Price     float64        `json:"price" validate:"required"`
	TraderId  enum.Trader    `json:"trader_id" validate:"required,len=24"`
}

func (s OrderData) Body() *OrderBody {
	return &OrderBody{
		Symbol:    s.Symbol,
		OrderType: s.OrderType,
		OrderSide: s.OrderSide,
		Amount:    s.Amount,
		Price:     s.Price,
	}
}

type OrderCancelBody struct {
	OrderId string `json:"order_id"`
	Symbol  string `json:"symbol"`
}

type OrderCancelData struct {
	Sandbox  bool        `json:"sandbox" validate:"omitempty"`
	Symbol   string      `json:"symbol" validate:"required"`
	TraderId enum.Trader `json:"trader_id" validate:"required,len=24"`
}

func (s OrderCancelData) Body() *OrderCancelBody {
	return &OrderCancelBody{
		Symbol: s.Symbol,
	}
}
