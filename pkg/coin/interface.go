package coin

import (
	"app/pkg/ccxt"
)

type ICoin interface {
	SetSandboxMode(enable any)
	SetMarkets(markets any, options ...any) any
	LoadMarkets(params ...any) <-chan any
	SignIn(options ...any) <-chan any
	GetSymbols() []string
	GetHas() map[string]any
	GetTimeframes() map[string]any
	GetCurrenciesList() []ccxt.Currency
	GetMarketsList() []ccxt.MarketInterface
	FetchTickers(options ...ccxt.FetchTickersOptions) (ccxt.Tickers, error)
	FetchOHLCV(symbol string, options ...ccxt.FetchOHLCVOptions) ([]ccxt.OHLCV, error)
	FetchTrades(symbol string, options ...ccxt.FetchTradesOptions) ([]ccxt.Trade, error)
	FetchOrderBook(symbol string, options ...ccxt.FetchOrderBookOptions) (ccxt.OrderBook, error)
	FetchBalance(params ...any) (ccxt.Balances, error)
	FetchClosedOrders(options ...ccxt.FetchClosedOrdersOptions) ([]ccxt.Order, error)
	FetchOpenOrders(options ...ccxt.FetchOpenOrdersOptions) ([]ccxt.Order, error)
	CreateOrder(symbol string, typeVar string, side string, amount float64, options ...ccxt.CreateOrderOptions) (ccxt.Order, error)
	CancelOrder(id string, options ...ccxt.CancelOrderOptions) (ccxt.Order, error)
}
