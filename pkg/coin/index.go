package coin

import (
	"app/pkg/ccxt"
	"app/pkg/ecode"
	"app/pkg/enum"
	"app/store"
	"context"
	"fmt"
	"maps"
	"slices"
	"strconv"

	"github.com/nhnghia272/gopkg"
)

type Config struct {
	Exchange enum.Exchange
	Tenant   enum.Tenant
	Trader   enum.Trader
	ApiKey   string
	Secret   string
	Password string
	Sandbox  bool
}

type TradeInfo struct {
	Exchange enum.Exchange
	Tenant   enum.Tenant
	Trader   enum.Trader
	Sandbox  bool
}

type Coins struct {
	store *store.Store
	items gopkg.CacheShard[*Coin]
}

func New(store *store.Store) *Coins {
	s := &Coins{store: store, items: gopkg.NewCacheShard[*Coin](64)}

	for _, exchange := range enum.ExchangeValues() {
		s.SetCoin(&Config{Exchange: exchange, Sandbox: true})
		s.SetCoin(&Config{Exchange: exchange, Sandbox: false})
	}

	return s
}

func (s *Coins) GetExchanges() []enum.Exchange {
	return enum.ExchangeValues()
}

func (s *Coins) SetCoin(cfg *Config) *Coin {
	coin := &Coin{
		ID:         cfg.Exchange,
		symbols:    make([]string, 0),
		timeframes: make([]string, 0),
		markets:    make([]*Market, 0),
		currencies: make([]*Currency, 0),
	}

	credential := make(map[string]any)
	if len(cfg.ApiKey) > 0 && len(cfg.Secret) > 0 {
		credential = map[string]any{"apiKey": cfg.ApiKey, "secret": cfg.Secret, "password": cfg.Password}
	}

	switch cfg.Exchange {
	case enum.ExchangeBinance:
		binance := ccxt.NewBinance(credential)
		coin.exchange = &binance
	case enum.ExchangeBitget:
		bitget := ccxt.NewBitget(credential)
		coin.exchange = &bitget
	case enum.ExchangeBybit:
		bybit := ccxt.NewBybit(credential)
		coin.exchange = &bybit
	case enum.ExchangeCex:
		cex := ccxt.NewCex(credential)
		coin.exchange = &cex
	case enum.ExchangeGate:
		gate := ccxt.NewGate(credential)
		coin.exchange = &gate
	case enum.ExchangeMexc:
		mexc := ccxt.NewMexc(credential)
		coin.exchange = &mexc
	case enum.ExchangeProbit:
		probit := ccxt.NewProbit(credential)
		coin.exchange = &probit
	}

	if item := coin.exchange.GetHas()["sandbox"]; item != nil && item.(bool) {
		coin.exchange.SetSandboxMode(cfg.Sandbox)
	}

	coin.LoadMarkets()

	s.items.Set(s.Key(&TradeInfo{Exchange: cfg.Exchange, Tenant: cfg.Tenant, Trader: cfg.Trader, Sandbox: cfg.Sandbox}), coin, -1)

	return coin
}

func (s *Coins) GetCoin(trade *TradeInfo) (*Coin, error) {
	coin, err := s.items.Get(s.Key(trade))
	if err != nil {
		trader, err := s.store.Db.Trader.FindOneById(context.Background(), string(trade.Trader))
		if err != nil {
			return nil, ecode.TraderNotFound
		}
		cfg := trader.ConfigDto()
		coin = s.SetCoin(&Config{Exchange: trade.Exchange, Tenant: trade.Tenant, Trader: trade.Trader, ApiKey: cfg.ApiKey, Secret: cfg.Secret, Password: cfg.Password, Sandbox: trade.Sandbox})
	}
	if item := coin.exchange.GetHas()["signIn"]; item != nil && item.(bool) {
		<-coin.exchange.SignIn()
	}
	return coin, nil
}

func (s *Coins) DelCoin(trade *TradeInfo) {
	s.items.Delete(s.Key(trade))
}

func (s *Coins) Key(trade *TradeInfo) string {
	return fmt.Sprintf("%v-%v-%v-%v", trade.Exchange, trade.Tenant, trade.Trader, trade.Sandbox)
}

type Coin struct {
	ID         enum.Exchange
	exchange   ICoin
	symbols    []string
	timeframes []string
	markets    []*Market
	currencies []*Currency
}

func (s *Coin) LoadMarkets() {
	<-s.exchange.LoadMarkets()
	s.symbols = slices.Sorted(slices.Values(s.exchange.GetSymbols()))
	s.timeframes = sortTimeframes(slices.Collect(maps.Keys(s.exchange.GetTimeframes())))
	s.markets = gopkg.MapFunc(s.exchange.GetMarketsList(), func(market ccxt.MarketInterface) *Market { return NewMarket(market) })
	s.currencies = gopkg.MapFunc(s.exchange.GetCurrenciesList(), func(currency ccxt.Currency) *Currency { return NewCurrency(currency) })
}

func (s *Coin) GetSymbols() []string {
	return s.symbols
}

func (s *Coin) GetTimeframes() []string {
	return s.timeframes
}

func (s *Coin) GetCurrencies() []*Currency {
	return s.currencies
}

func (s *Coin) GetMarkets(query *MarketQuery) []*Market {
	return gopkg.FilterFunc(s.markets, func(market *Market) bool {
		if query.MarketType == enum.MarketTypeSpot {
			return market.Spot
		}
		if query.MarketType == enum.MarketTypeMargin {
			return market.Margin
		}
		return false
	})
}

func (s *Coin) GetTickers(query *TickerQuery) ([]*Ticker, error) {
	res, err := s.exchange.FetchTickers(ccxt.WithFetchTickersSymbols(query.Symbols))
	if err != nil {
		return nil, ecode.BadRequest.Desc(err)
	}
	return gopkg.MapFunc(slices.Collect(maps.Values(res.Tickers)), func(ticker ccxt.Ticker) *Ticker { return NewTicker(ticker) }), nil
}

func (s *Coin) GetOHLCVs(query *OHLCVQuery) ([]*OHLCV, error) {
	res, err := s.exchange.FetchOHLCV(query.Symbol, ccxt.WithFetchOHLCVTimeframe(query.Timeframe))
	if err != nil {
		return nil, ecode.BadRequest.Desc(err)
	}
	return gopkg.MapFunc(res, func(ohlcv ccxt.OHLCV) *OHLCV { return NewOHLCV(ohlcv) }), nil
}

func (s *Coin) GetTrades(query *TradeQuery) ([]*Trade, error) {
	res, err := s.exchange.FetchTrades(query.Symbol)
	if err != nil {
		return nil, ecode.BadRequest.Desc(err)
	}
	return gopkg.MapFunc(res, func(trade ccxt.Trade) *Trade { return NewTrade(trade) }), nil
}

func (s *Coin) GetOrderBook(query *OrderBookQuery) (*OrderBook, error) {
	res, err := s.exchange.FetchOrderBook(query.Symbol)
	if err != nil {
		return nil, ecode.BadRequest.Desc(err)
	}
	return NewOrderBook(res), nil
}

func (s *Coin) GetBalances(query *BalanceQuery) ([]*Balance, error) {
	res, err := s.exchange.FetchBalance(map[string]any{"type": query.MarketType})
	if err != nil {
		return nil, ecode.BadRequest.Desc(err)
	}
	return gopkg.MapFunc(slices.Collect(maps.Keys(res.Balances)), func(key string) *Balance { return NewBalance(key, res.Balances[key]) }), nil
}

func (s *Coin) GetOrders(query *OrderQuery) ([]*Order, error) {
	end := ccxt.Milliseconds()
	start := end - 1000*60*60*24*90

	optsClosed := func(opts *ccxt.FetchClosedOrdersOptionsStruct) {
		opts.Symbol = gopkg.Pointer(query.Symbol)
		if s.ID == enum.ExchangeProbit {
			opts.Params = gopkg.Pointer(map[string]any{"start_time": ccxt.Iso8601(start), "end_time": ccxt.Iso8601(end), "limit": 100})
		}
	}

	ordersClosed, err := s.exchange.FetchClosedOrders(optsClosed)
	if err != nil {
		return nil, ecode.BadRequest.Desc(err)
	}

	optsOpen := func(opts *ccxt.FetchOpenOrdersOptionsStruct) {
		opts.Symbol = gopkg.Pointer(query.Symbol)
		if s.ID == enum.ExchangeProbit {
			opts.Params = gopkg.Pointer(map[string]any{"start_time": ccxt.Iso8601(start), "end_time": ccxt.Iso8601(end), "limit": 100})
		}
	}

	orders, err := s.exchange.FetchOpenOrders(optsOpen)
	if err != nil {
		return nil, ecode.BadRequest.Desc(err)
	}

	orders = append(orders, ordersClosed...)

	slices.SortFunc(orders, func(a, b ccxt.Order) int {
		idA, _ := strconv.Atoi(gopkg.Value(a.Id))
		idB, _ := strconv.Atoi(gopkg.Value(b.Id))
		return idB - idA
	})

	return gopkg.MapFunc(orders, func(order ccxt.Order) *Order { return NewOrder(order) }), nil
}

func (s *Coin) CreateOrder(body *OrderBody) (*Order, error) {
	res, err := s.exchange.CreateOrder(body.Symbol, string(body.OrderType), string(body.OrderSide), body.Amount, ccxt.WithCreateOrderPrice(body.Price))
	if err != nil {
		return nil, ecode.BadRequest.Desc(err)
	}
	return NewOrder(res), nil
}

func (s *Coin) CancelOrder(body *OrderCancelBody) (*Order, error) {
	res, err := s.exchange.CancelOrder(body.OrderId, ccxt.WithCancelOrderSymbol(body.Symbol))
	if err != nil {
		return nil, ecode.BadRequest.Desc(err)
	}
	return NewOrder(res), nil
}
