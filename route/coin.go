package route

import (
	"app/env"
	"app/pkg/api"
	"app/pkg/coin"
	"app/pkg/ecode"
	"app/pkg/enum"
	"app/pkg/tz"
	"app/store/db"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/leekchan/accounting"
	"github.com/nhnghia272/gopkg"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

type coins struct {
	*middleware
}

func init() {
	handlers = append(handlers, func(m *middleware, r *gin.Engine) {
		s := coins{m}
		if len(env.SlackUri) > 0 {
			s.autoAlertOrder()
		}

		v1cms := r.Group("/cms/v1/coins")
		v1cms.GET("/configs", s.BearerAuth(enum.PermissionCoinView), s.v1cms_ListConfig())
		v1cms.POST("/configs", s.BearerAuth(enum.PermissionCoinCreate), s.v1cms_CreateConfig())
		v1cms.PUT("/configs/:config_id", s.BearerAuth(enum.PermissionCoinUpdate), s.v1cms_UpdateConfig())
		v1cms.GET("/orders", s.BearerAuth(enum.PermissionCoinView), s.v1cms_ListOrder())
		v1cms.GET("/reports", s.BearerAuth(enum.PermissionCoinView), s.v1cms_ListReport())
		v1cms.GET("/traders", s.BearerAuth(enum.PermissionCoinView), s.v1cms_ListTrader())
		v1cms.GET("/strategies", s.BearerAuth(enum.PermissionCoinView), s.v1cms_ListStrategy())
		v1cms.GET("/exchanges", s.BearerAuth(enum.PermissionCoinView), s.v1cms_ListExchange())
		v1cms.GET("/exchanges/:exchange_id/symbols", s.BearerAuth(enum.PermissionCoinView), s.v1cms_ListSymbolByExchange())
		v1cms.GET("/exchanges/:exchange_id/timeframes", s.BearerAuth(enum.PermissionCoinView), s.v1cms_ListTimeframeByExchange())
		v1cms.GET("/exchanges/:exchange_id/currencies", s.BearerAuth(enum.PermissionCoinView), s.v1cms_ListCurrencyByExchange())
		v1cms.GET("/exchanges/:exchange_id/markets", s.BearerAuth(enum.PermissionCoinView), s.v1cms_ListMarketByExchange())
		v1cms.GET("/exchanges/:exchange_id/tickers", s.BearerAuth(enum.PermissionCoinView), s.v1cms_ListTickerByExchange())
		v1cms.GET("/exchanges/:exchange_id/ohlcvs", s.BearerAuth(enum.PermissionCoinView), s.v1cms_ListOHLCVByExchange())
		v1cms.GET("/exchanges/:exchange_id/trades", s.BearerAuth(enum.PermissionCoinView), s.v1cms_ListTradeByExchange())
		v1cms.GET("/exchanges/:exchange_id/orderbook", s.BearerAuth(enum.PermissionCoinView), s.v1cms_OrderBookByExchange())
		v1cms.GET("/exchanges/:exchange_id/balances", s.BearerAuth(enum.PermissionCoinView), s.v1cms_ListBalanceByExchange())
		v1cms.GET("/exchanges/:exchange_id/orders", s.BearerAuth(enum.PermissionCoinView), s.v1cms_ListOrderByExchange())

		v1agent := r.Group("/agent/v1/coins")
		v1agent.GET("/configs", s.BasicAuth(), s.v1agent_ListConfig())
		v1agent.PUT("/configs/:config_id/server-status", s.BasicAuth(), s.v1agent_UpdateServerStatus())
		v1agent.POST("/orders", s.BasicAuth(), s.v1agent_CreateOrder())
	})
}

// @Tags Cms
// @Summary List Configs
// @Security BearerAuth
// @Param query query db.CoinConfigCmsQuery false "query"
// @Success 200 {object} []db.CoinConfigCmsDto
// @Failure 400,401,403,500 {object} ecode.Error
// @Router /cms/v1/coins/configs [get]
func (s coins) v1cms_ListConfig() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := s.Session(c)

		qb := &db.CoinConfigCmsQuery{}
		if err := c.ShouldBind(qb); err != nil {
			c.Error(ecode.BadRequest.Desc(err))
			return
		}

		query := qb.BuildCore()
		query.TenantId = gopkg.Pointer(session.TenantId)

		domains, _ := s.store.Db.CoinConfig.FindAll(c.Request.Context(), query)
		results := gopkg.MapFunc(domains, func(domain *db.CoinConfigDomain) *db.CoinConfigCmsDto { return domain.CmsDto() })

		s.Pagination(c, s.store.Db.CoinConfig.Count(c.Request.Context(), query), query.Page, query.Limit)

		c.JSON(http.StatusOK, results)
	}
}

// @Tags Cms
// @Summary Create Config
// @Security BearerAuth
// @Param body body db.CoinConfigCmsData true "body"
// @Success 200 {object} db.CoinConfigCmsDto
// @Failure 400,401,403,409,500 {object} ecode.Error
// @Router /cms/v1/coins/configs [post]
func (s coins) v1cms_CreateConfig() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := s.Session(c)

		data := &db.CoinConfigCmsData{}
		if err := c.ShouldBind(data); err != nil {
			c.Error(ecode.BadRequest.Desc(err))
			return
		}

		strategy, err := s.store.Db.Strategy.FindOneById(c.Request.Context(), data.StrategyId)
		if err != nil || gopkg.Value(strategy.DataStatus) == enum.DataStatusDisable {
			c.Error(ecode.StrategyDisabled)
			return
		}

		for _, id := range data.TraderIds {
			trader, err := s.store.Db.Trader.FindOneById(c.Request.Context(), id)
			if err != nil || gopkg.Value(trader.DataStatus) == enum.DataStatusDisable {
				c.Error(ecode.TraderDisabled)
				return
			}

			if gopkg.Value(trader.TenantId) != session.TenantId || gopkg.Value(trader.Exchange) != data.Exchange || gopkg.Value(trader.Sandbox) != data.Sandbox {
				c.Error(ecode.TraderNotFound)
				return
			}
		}

		domain := data.Domain(&db.CoinConfigDomain{})
		domain.TenantId = gopkg.Pointer(session.TenantId)
		domain.CreatedBy = gopkg.Pointer(session.Username)
		domain.UpdatedBy = gopkg.Pointer(session.Username)

		domain, err = s.store.Db.CoinConfig.Save(c.Request.Context(), domain)
		if err != nil {
			c.Error(ecode.ConfigConflict.Stack(err))
			return
		}

		s.AuditLog(c, s.store.Db.CoinConfig.CollectionName(), enum.DataActionCreate, data, domain, db.SID(domain.ID))

		s.store.RunStrategy(c.Request.Context(), domain)

		c.JSON(http.StatusOK, domain.CmsDto())
	}
}

// @Tags Cms
// @Summary Update Config
// @Security BearerAuth
// @Param config_id path string true "config_id"
// @Param body body db.CoinConfigCmsData true "body"
// @Success 200 {object} db.CoinConfigCmsDto
// @Failure 400,401,403,404,409,500 {object} ecode.Error
// @Router /cms/v1/coins/configs/{config_id} [put]
func (s coins) v1cms_UpdateConfig() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := s.Session(c)

		data := &db.CoinConfigCmsData{}
		if err := c.ShouldBind(data); err != nil {
			c.Error(ecode.BadRequest.Desc(err))
			return
		}

		strategy, err := s.store.Db.Strategy.FindOneById(c.Request.Context(), data.StrategyId)
		if err != nil || gopkg.Value(strategy.DataStatus) == enum.DataStatusDisable {
			c.Error(ecode.StrategyDisabled)
			return
		}

		for _, id := range data.TraderIds {
			trader, err := s.store.Db.Trader.FindOneById(c.Request.Context(), id)
			if err != nil || gopkg.Value(trader.DataStatus) == enum.DataStatusDisable {
				c.Error(ecode.TraderDisabled)
				return
			}

			if gopkg.Value(trader.TenantId) != session.TenantId || gopkg.Value(trader.Exchange) != data.Exchange || gopkg.Value(trader.Sandbox) != data.Sandbox {
				c.Error(ecode.TraderNotFound)
				return
			}
		}

		domain, err := s.store.Db.CoinConfig.FindOneById(c.Request.Context(), c.Param("config_id"))
		if err != nil || gopkg.Value(domain.TenantId) != session.TenantId {
			c.Error(ecode.ConfigNotFound)
			return
		}

		update := &db.CoinConfigDomain{BaseDomain: db.BaseDomain{ID: domain.ID, UpdatedBy: gopkg.Pointer(session.Username)}}
		update = data.Domain(update)

		domain, err = s.store.Db.CoinConfig.Save(c.Request.Context(), update)
		if err != nil {
			c.Error(ecode.ConfigConflict.Stack(err))
			return
		}

		s.AuditLog(c, s.store.Db.CoinConfig.CollectionName(), enum.DataActionUpdate, data, domain, db.SID(domain.ID))

		s.store.RunStrategy(c.Request.Context(), domain)

		c.JSON(http.StatusOK, domain.CmsDto())
	}
}

// @Tags Cms
// @Summary List Orders
// @Security BearerAuth
// @Param query query db.CoinOrderCmsQuery false "query"
// @Success 200 {object} []db.CoinOrderCmsDto
// @Router /cms/v1/coins/orders [get]
func (s coins) v1cms_ListOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := s.Session(c)

		qb := &db.CoinOrderCmsQuery{}
		if err := c.ShouldBind(qb); err != nil {
			c.Error(ecode.BadRequest.Desc(err))
			return
		}

		query := qb.BuildCore(session.Timezone)
		query.TenantId = gopkg.Pointer(session.TenantId)

		domains, _ := s.store.Db.CoinOrder.FindAll(c.Request.Context(), query)
		results := gopkg.MapFunc(domains, func(domain *db.CoinOrderDomain) *db.CoinOrderCmsDto { return domain.CmsDto() })

		s.Pagination(c, s.store.Db.CoinOrder.Count(c.Request.Context(), query), query.Page, query.Limit)

		c.JSON(http.StatusOK, results)
	}
}

// @Tags Cms
// @Summary List Reports
// @Security BearerAuth
// @Param query query db.CoinOrderCmsReportQuery false "query"
// @Success 200 {object} []db.CoinOrderCmsReportDto
// @Router /cms/v1/coins/reports [get]
func (s coins) v1cms_ListReport() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := s.Session(c)

		qb := &db.CoinOrderCmsReportQuery{}
		if err := c.ShouldBind(qb); err != nil {
			c.Error(ecode.BadRequest.Desc(err))
			return
		}

		config, err := s.store.Db.CoinConfig.FindOneById(c.Request.Context(), gopkg.Value(qb.ConfigId))
		if err != nil || gopkg.Value(config.TenantId) != session.TenantId {
			c.Error(ecode.ConfigNotFound)
			return
		}

		query := qb.BuildCore(session.Timezone)
		query.TenantId = gopkg.Pointer(session.TenantId)
		query.Status = gopkg.Pointer(enum.OrderStatusClosed)

		report, err := s.store.Db.CoinOrder.ReportCoinOrderDetails(c.Request.Context(), query)
		if err != nil {
			c.Error(ecode.InternalServerError.Stack(err))
			return
		}

		mycoin, err := s.coin.GetCoin(&coin.TradeInfo{Exchange: enum.Exchange(gopkg.Value(config.Exchange)), Sandbox: gopkg.Value(config.Sandbox)})
		if err != nil {
			c.Error(ecode.ExchangeNotFound)
			return
		}

		tickers, err := mycoin.GetTickers(&coin.TickerQuery{Sandbox: gopkg.Value(config.Sandbox), Symbols: gopkg.Value(config.Symbols)})
		if err != nil {
			c.Error(ecode.BadRequest.Desc(err))
			return
		}

		for _, ticker := range tickers {
			report.Summary24H.TotalTokenExchange += ticker.BaseVolume
			report.Summary24H.TotalVolumeExchange += ticker.QuoteVolume
		}

		c.JSON(http.StatusOK, report)
	}
}

// @Tags Cms
// @Summary List Traders
// @Security BearerAuth
// @Success 200 {object} []db.TraderBaseDto
// @Router /cms/v1/coins/traders [get]
func (s coins) v1cms_ListTrader() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := s.Session(c)

		query := &db.TraderQuery{Query: db.Query{Sorts: "name:asc"}}
		query.TenantId = gopkg.Pointer(session.TenantId)

		domains, _ := s.store.Db.Trader.FindAll(c.Request.Context(), query)
		results := gopkg.MapFunc(domains, func(domain *db.TraderDomain) *db.TraderBaseDto { return domain.BaseDto() })

		c.JSON(http.StatusOK, results)
	}
}

// @Tags Cms
// @Summary List Strategies
// @Security BearerAuth
// @Success 200 {object} []db.StrategyCmsDto
// @Router /cms/v1/coins/strategies [get]
func (s coins) v1cms_ListStrategy() gin.HandlerFunc {
	return func(c *gin.Context) {
		query := &db.StrategyQuery{Query: db.Query{Sorts: "name:asc"}}

		domains, _ := s.store.Db.Strategy.FindAll(c.Request.Context(), query)
		results := gopkg.MapFunc(domains, func(domain *db.StrategyDomain) *db.StrategyCmsDto { return domain.CmsDto() })

		c.JSON(http.StatusOK, results)
	}
}

// @Tags Cms
// @Summary List Exchanges
// @Security BearerAuth
// @Success 200 {object} []enum.Exchange
// @Router /cms/v1/coins/exchanges [get]
func (s coins) v1cms_ListExchange() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, s.coin.GetExchanges())
	}
}

// @Tags Cms
// @Summary List Symbols
// @Security BearerAuth
// @Param exchange_id path enum.Exchange true "exchange_id"
// @Param query query coin.SymbolQuery false "query"
// @Success 200 {object} []string
// @Router /cms/v1/coins/exchanges/{exchange_id}/symbols [get]
func (s coins) v1cms_ListSymbolByExchange() gin.HandlerFunc {
	return func(c *gin.Context) {
		query := &coin.SymbolQuery{}
		if err := c.ShouldBind(query); err != nil {
			c.Error(ecode.BadRequest.Desc(err))
			return
		}

		mycoin, err := s.coin.GetCoin(&coin.TradeInfo{Exchange: enum.Exchange(c.Param("exchange_id")), Sandbox: query.Sandbox})
		if err != nil {
			c.Error(ecode.ExchangeNotFound)
			return
		}

		c.JSON(http.StatusOK, mycoin.GetSymbols())
	}
}

// @Tags Cms
// @Summary List Timeframes
// @Security BearerAuth
// @Param exchange_id path enum.Exchange true "exchange_id"
// @Param query query coin.TimeframeQuery false "query"
// @Success 200 {object} []string
// @Router /cms/v1/coins/exchanges/{exchange_id}/timeframes [get]
func (s coins) v1cms_ListTimeframeByExchange() gin.HandlerFunc {
	return func(c *gin.Context) {
		query := &coin.TimeframeQuery{}
		if err := c.ShouldBind(query); err != nil {
			c.Error(ecode.BadRequest.Desc(err))
			return
		}

		mycoin, err := s.coin.GetCoin(&coin.TradeInfo{Exchange: enum.Exchange(c.Param("exchange_id")), Sandbox: query.Sandbox})
		if err != nil {
			c.Error(ecode.ExchangeNotFound)
			return
		}

		c.JSON(http.StatusOK, mycoin.GetTimeframes())
	}
}

// @Tags Cms
// @Summary List Currencies
// @Security BearerAuth
// @Param exchange_id path enum.Exchange true "exchange_id"
// @Param query query coin.CurrencyQuery false "query"
// @Success 200 {object} []coin.Currency
// @Router /cms/v1/coins/exchanges/{exchange_id}/currencies [get]
func (s coins) v1cms_ListCurrencyByExchange() gin.HandlerFunc {
	return func(c *gin.Context) {
		query := &coin.CurrencyQuery{}
		if err := c.ShouldBind(query); err != nil {
			c.Error(ecode.BadRequest.Desc(err))
			return
		}

		mycoin, err := s.coin.GetCoin(&coin.TradeInfo{Exchange: enum.Exchange(c.Param("exchange_id")), Sandbox: query.Sandbox})
		if err != nil {
			c.Error(ecode.ExchangeNotFound)
			return
		}

		c.JSON(http.StatusOK, mycoin.GetCurrencies())
	}
}

// @Tags Cms
// @Summary List Markets
// @Security BearerAuth
// @Param exchange_id path enum.Exchange true "exchange_id"
// @Param query query coin.MarketQuery false "query"
// @Success 200 {object} []coin.Market
// @Router /cms/v1/coins/exchanges/{exchange_id}/markets [get]
func (s coins) v1cms_ListMarketByExchange() gin.HandlerFunc {
	return func(c *gin.Context) {
		query := &coin.MarketQuery{}
		if err := c.ShouldBind(query); err != nil {
			c.Error(ecode.BadRequest.Desc(err))
			return
		}

		mycoin, err := s.coin.GetCoin(&coin.TradeInfo{Exchange: enum.Exchange(c.Param("exchange_id")), Sandbox: query.Sandbox})
		if err != nil {
			c.Error(ecode.ExchangeNotFound)
			return
		}

		c.JSON(http.StatusOK, mycoin.GetMarkets(query))
	}
}

// @Tags Cms
// @Summary List Tickers
// @Security BearerAuth
// @Param exchange_id path enum.Exchange true "exchange_id"
// @Param query query coin.TickerQuery false "query"
// @Success 200 {object} []coin.Ticker
// @Router /cms/v1/coins/exchanges/{exchange_id}/tickers [get]
func (s coins) v1cms_ListTickerByExchange() gin.HandlerFunc {
	return func(c *gin.Context) {
		query := &coin.TickerQuery{}
		if err := c.ShouldBind(query); err != nil {
			c.Error(ecode.BadRequest.Desc(err))
			return
		}

		mycoin, err := s.coin.GetCoin(&coin.TradeInfo{Exchange: enum.Exchange(c.Param("exchange_id")), Sandbox: query.Sandbox})
		if err != nil {
			c.Error(ecode.ExchangeNotFound)
			return
		}

		tickers, err := mycoin.GetTickers(query)
		if err != nil {
			c.Error(ecode.BadRequest.Desc(err))
			return
		}

		c.JSON(http.StatusOK, tickers)
	}
}

// @Tags Cms
// @Summary List OHLCVs
// @Security BearerAuth
// @Param exchange_id path enum.Exchange true "exchange_id"
// @Param query query coin.OHLCVQuery false "query"
// @Success 200 {object} []coin.OHLCV
// @Router /cms/v1/coins/exchanges/{exchange_id}/ohlcvs [get]
func (s coins) v1cms_ListOHLCVByExchange() gin.HandlerFunc {
	return func(c *gin.Context) {
		query := &coin.OHLCVQuery{}
		if err := c.ShouldBind(query); err != nil {
			c.Error(ecode.BadRequest.Desc(err))
		}

		mycoin, err := s.coin.GetCoin(&coin.TradeInfo{Exchange: enum.Exchange(c.Param("exchange_id")), Sandbox: query.Sandbox})
		if err != nil {
			c.Error(ecode.ExchangeNotFound)
			return
		}

		ohlcvs, err := mycoin.GetOHLCVs(query)
		if err != nil {
			c.Error(ecode.BadRequest.Desc(err))
			return
		}

		c.JSON(http.StatusOK, ohlcvs)
	}
}

// @Tags Cms
// @Summary List Trades
// @Security BearerAuth
// @Param exchange_id path enum.Exchange true "exchange_id"
// @Param query query coin.TradeQuery false "query"
// @Success 200 {object} []coin.Trade
// @Router /cms/v1/coins/exchanges/{exchange_id}/trades [get]
func (s coins) v1cms_ListTradeByExchange() gin.HandlerFunc {
	return func(c *gin.Context) {
		query := &coin.TradeQuery{}
		if err := c.ShouldBind(query); err != nil {
			c.Error(ecode.BadRequest.Desc(err))
			return
		}

		mycoin, err := s.coin.GetCoin(&coin.TradeInfo{Exchange: enum.Exchange(c.Param("exchange_id")), Sandbox: query.Sandbox})
		if err != nil {
			c.Error(ecode.ExchangeNotFound)
			return
		}

		trades, err := mycoin.GetTrades(query)
		if err != nil {
			c.Error(ecode.BadRequest.Desc(err))
			return
		}

		c.JSON(http.StatusOK, trades)
	}
}

// @Tags Cms
// @Summary Get OrderBook
// @Security BearerAuth
// @Param exchange_id path enum.Exchange true "exchange_id"
// @Param query query coin.OrderBookQuery false "query"
// @Success 200 {object} coin.OrderBook
// @Router /cms/v1/coins/exchanges/{exchange_id}/orderbook [get]
func (s coins) v1cms_OrderBookByExchange() gin.HandlerFunc {
	return func(c *gin.Context) {
		query := &coin.OrderBookQuery{}
		if err := c.ShouldBind(query); err != nil {
			c.Error(ecode.BadRequest.Desc(err))
			return
		}

		mycoin, err := s.coin.GetCoin(&coin.TradeInfo{Exchange: enum.Exchange(c.Param("exchange_id")), Sandbox: query.Sandbox})
		if err != nil {
			c.Error(ecode.ExchangeNotFound)
			return
		}

		orderbook, err := mycoin.GetOrderBook(query)
		if err != nil {
			c.Error(ecode.BadRequest.Desc(err))
			return
		}

		c.JSON(http.StatusOK, orderbook)
	}
}

// @Tags Cms
// @Summary List Balances
// @Security BearerAuth
// @Param exchange_id path enum.Exchange true "exchange_id"
// @Param query query coin.BalanceQuery false "query"
// @Success 200 {object} []coin.Balance
// @Router /cms/v1/coins/exchanges/{exchange_id}/balances [get]
func (s coins) v1cms_ListBalanceByExchange() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := s.Session(c)

		query := &coin.BalanceQuery{}
		if err := c.ShouldBind(query); err != nil {
			c.Error(ecode.BadRequest.Desc(err))
			return
		}

		mycoin, err := s.coin.GetCoin(&coin.TradeInfo{Exchange: enum.Exchange(c.Param("exchange_id")), Sandbox: query.Sandbox, Trader: query.TraderId, Tenant: session.TenantId})
		if err != nil {
			c.Error(ecode.ExchangeNotFound)
			return
		}

		balances, err := mycoin.GetBalances(query)
		if err != nil {
			c.Error(ecode.BadRequest.Desc(err))
			return
		}

		c.JSON(http.StatusOK, balances)
	}
}

// @Tags Cms
// @Summary List Orders
// @Security BearerAuth
// @Param exchange_id path enum.Exchange true "exchange_id"
// @Param query query coin.OrderQuery false "query"
// @Success 200 {object} []coin.Order
// @Router /cms/v1/coins/exchanges/{exchange_id}/orders [get]
func (s coins) v1cms_ListOrderByExchange() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := s.Session(c)

		query := &coin.OrderQuery{}
		if err := c.ShouldBind(query); err != nil {
			c.Error(ecode.BadRequest.Desc(err))
			return
		}

		mycoin, err := s.coin.GetCoin(&coin.TradeInfo{Exchange: enum.Exchange(c.Param("exchange_id")), Sandbox: query.Sandbox, Trader: query.TraderId, Tenant: session.TenantId})
		if err != nil {
			c.Error(ecode.ExchangeNotFound)
			return
		}

		orders, err := mycoin.GetOrders(query)
		if err != nil {
			c.Error(ecode.BadRequest.Desc(err))
			return
		}

		c.JSON(http.StatusOK, orders)
	}
}

// @Tags Agent
// @Summary List Configs
// @Security BasicAuth
// @Success 200 {object} []db.CoinConfigStrategyDto
// @Router /agent/v1/coins/configs [get]
func (s coins) v1agent_ListConfig() gin.HandlerFunc {
	return func(c *gin.Context) {
		domains, _ := s.store.Db.CoinConfig.FindAll(c.Request.Context(), &db.CoinConfigQuery{DataStatus: gopkg.Pointer(enum.DataStatusEnable)})
		results, _ := gopkg.MapParallelFunc(domains, func(domain *db.CoinConfigDomain) *db.CoinConfigStrategyDto {
			var (
				strategy, _ = s.store.Db.Strategy.FindOneById(c.Request.Context(), gopkg.Value(domain.StrategyId))
				strategyDto = domain.StrategyDto()
			)

			strategyDto.Attributes["filename"] = gopkg.Value(strategy.Filename)

			gopkg.LoopFunc(gopkg.Value(domain.TraderIds), func(id string) {
				trader, _ := s.store.Db.Trader.FindOneById(c.Request.Context(), id)
				traderDto := trader.ConfigDto()

				strategyDto.Credentials = append(strategyDto.Credentials, db.CoinConfigCredentialInfo{
					TraderId: traderDto.TraderId,
					ApiKey:   traderDto.ApiKey,
					Secret:   traderDto.Secret,
					Password: traderDto.Password,
				})
			})

			return strategyDto
		})
		c.JSON(http.StatusOK, results)
	}
}

// @Tags Agent
// @Summary Update Server Status
// @Security BasicAuth
// @Param config_id path string true "config_id"
// @Param body body db.CoinConfigServerInfo true "body"
// @Success 200
// @Router /agent/v1/coins/configs/{config_id}/server-status [put]
func (s coins) v1agent_UpdateServerStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		data := &db.CoinConfigServerInfo{}
		if err := c.ShouldBind(data); err != nil {
			c.Error(ecode.BadRequest.Desc(err))
			return
		}

		domain, err := s.store.Db.CoinConfig.FindOneById(c.Request.Context(), c.Param("config_id"))
		if err != nil {
			c.Error(ecode.ConfigNotFound)
			return
		}

		update := &db.CoinConfigDomain{BaseDomain: db.BaseDomain{ID: domain.ID}}
		update.Server = data

		if _, err = s.store.Db.CoinConfig.Save(c.Request.Context(), update); err != nil {
			c.Error(ecode.ConfigConflict.Stack(err))
			return
		}

		c.Status(http.StatusOK)
	}
}

// @Tags Agent
// @Summary Create Order
// @Security BasicAuth
// @Param body body db.CoinOrderAlertData true "body"
// @Success 200
// @Router /agent/v1/coins/orders [post]
func (s coins) v1agent_CreateOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := s.Session(c)

		data := &db.CoinOrderAlertData{}
		if err := c.ShouldBind(data); err != nil {
			c.Error(ecode.BadRequest.Desc(err))
			return
		}

		domain := data.Domain(&db.CoinOrderDomain{})
		domain.CreatedBy = gopkg.Pointer(session.Username)
		domain.UpdatedBy = gopkg.Pointer(session.Username)

		if _, err := s.store.Db.CoinOrder.Save(c.Request.Context(), domain); err != nil {
			c.Error(ecode.BadRequest.Desc(err))
			return
		}

		c.Status(http.StatusOK)
	}
}

func (s *coins) autoAlertOrder() {
	var (
		timezone = tz.AsiaHoChiMinh
		location = tz.Location(timezone)
	)

	c := cron.New(cron.WithLocation(location))
	c.AddFunc("0 8,14,22 * * *", func() {
		var (
			ctx        = context.Background()
			configs, _ = s.store.Db.CoinConfig.FindAll(ctx, &db.CoinConfigQuery{DataStatus: gopkg.Pointer(enum.DataStatusEnable)})
		)
		for _, config := range configs {
			qb := &db.CoinOrderCmsReportQuery{}
			qb.ConfigId = gopkg.Pointer(db.SID(config.ID))
			qb.EndDate = gopkg.Pointer(time.Now().In(location).Format(time.DateOnly))
			qb.StartDate = gopkg.Pointer(time.Now().In(location).Format(time.DateOnly))

			query := qb.BuildCore(timezone)
			query.TenantId = config.TenantId
			query.Status = gopkg.Pointer(enum.OrderStatusClosed)

			report, err := s.store.Db.CoinOrder.ReportCoinOrderDetails(ctx, query)
			if err != nil {
				logrus.Errorln("autoAlertOrder", err)
				return
			}

			mycoin, err := s.coin.GetCoin(&coin.TradeInfo{Exchange: enum.Exchange(gopkg.Value(config.Exchange)), Sandbox: gopkg.Value(config.Sandbox)})
			if err != nil {
				logrus.Errorln("autoAlertOrder", err)
				return
			}

			tickers, err := mycoin.GetTickers(&coin.TickerQuery{Sandbox: gopkg.Value(config.Sandbox), Symbols: gopkg.Value(config.Symbols)})
			if err != nil {
				logrus.Errorln("autoAlertOrder", err)
				return
			}

			for _, ticker := range tickers {
				report.Summary24H.TotalTokenExchange += ticker.BaseVolume
				report.Summary24H.TotalVolumeExchange += ticker.QuoteVolume
			}

			money := accounting.Accounting{Precision: 3}

			text := fmt.Sprintf("```|---------------------------|\n"+
				"| %-25s |\n"+
				"|---------------------------|\n"+
				"| %-6s | %16s |\n"+
				"| %-6s | %16s |\n"+
				"| %-6s | %16s |\n"+
				"|---------------------------|\n"+
				"| %-6s | %16s |\n"+
				"| %-6s | %16s |\n"+
				"| %-6s | %16s |\n"+
				"|---------------------------|\n"+
				"| %-6s | %16s |\n"+
				"| %-6s | %16s |\n"+
				"| %-6s | %16s |\n"+
				"| %-6s | %16s |\n"+
				"| %-6s | %16.0f |\n"+
				"|---------------------------|\n```",
				strings.ToUpper(string(gopkg.Value(config.Exchange))),
				"Name", "Exchange",
				"Volume", money.FormatMoney(report.Summary24H.TotalVolumeExchange),
				"Token", money.FormatMoney(report.Summary24H.TotalTokenExchange),
				"Name", "Others",
				"Volume", money.FormatMoney(report.Summary24H.TotalVolumeExchange-report.Summary24H.TotalVolume),
				"Token", money.FormatMoney(report.Summary24H.TotalTokenExchange-report.Summary24H.TotalToken),
				"Name", "Vertix",
				"Volume", money.FormatMoney(report.Summary24H.TotalVolume),
				"Token", money.FormatMoney(report.Summary24H.TotalToken),
				"Fee", money.FormatMoney(report.Summary24H.TotalFee),
				"Order", report.Summary24H.TotalOrder,
			)

			body, _ := json.Marshal(map[string]string{"text": string(text)})

			api.New[any](http.NewRequest(http.MethodPost, env.SlackUri, bytes.NewBuffer(body))).Call()
		}
	})
	c.Start()
}
