package enum

import (
	"slices"

	"github.com/nhnghia272/gopkg"
)

// Tenant
type Tenant string

// Trader
type Trader string

// Strategy
type Strategy string

// Kind
type Kind string

const (
	KindPermission   Kind = "permission"
	KindDataStatus   Kind = "data_status"
	KindDataAction   Kind = "data_action"
	KindExchange     Kind = "exchange"
	KindMarketType   Kind = "market_type"
	KindMarginMode   Kind = "margin_mode"
	KindOrderType    Kind = "order_type"
	KindOrderSide    Kind = "order_side"
	KindOrderStatus  Kind = "order_status"
	KindServerStatus Kind = "server_status"
)

func Tags() map[string][]string {
	return map[string][]string{
		string(KindPermission):   gopkg.MapFunc(PermissionValues(), func(e Permission) string { return string(e) }),
		string(KindDataStatus):   gopkg.MapFunc(DataStatusValues(), func(e DataStatus) string { return string(e) }),
		string(KindDataAction):   gopkg.MapFunc(DataActionValues(), func(e DataAction) string { return string(e) }),
		string(KindExchange):     gopkg.MapFunc(ExchangeValues(), func(e Exchange) string { return string(e) }),
		string(KindMarketType):   gopkg.MapFunc(MarketTypeValues(), func(e MarketType) string { return string(e) }),
		string(KindMarginMode):   gopkg.MapFunc(MarginModeValues(), func(e MarginMode) string { return string(e) }),
		string(KindOrderType):    gopkg.MapFunc(OrderTypeValues(), func(e OrderType) string { return string(e) }),
		string(KindOrderSide):    gopkg.MapFunc(OrderSideValues(), func(e OrderSide) string { return string(e) }),
		string(KindOrderStatus):  gopkg.MapFunc(OrderStatusValues(), func(e OrderStatus) string { return string(e) }),
		string(KindServerStatus): gopkg.MapFunc(ServerStatusValues(), func(e ServerStatus) string { return string(e) }),
	}
}

// Permission
type Permission string

const (
	PermissionSystemSetting  Permission = "system_setting"
	PermissionSystemAuditLog Permission = "system_audit_log"

	PermissionClientView   Permission = "client_view"
	PermissionClientCreate Permission = "client_create"
	PermissionClientDelete Permission = "client_delete"

	PermissionRoleView   Permission = "role_view"
	PermissionRoleCreate Permission = "role_create"
	PermissionRoleUpdate Permission = "role_update"

	PermissionUserView   Permission = "user_view"
	PermissionUserCreate Permission = "user_create"
	PermissionUserUpdate Permission = "user_update"

	PermissionTenantView   Permission = "tenant_view"
	PermissionTenantCreate Permission = "tenant_create"
	PermissionTenantUpdate Permission = "tenant_update"

	PermissionCoinView   Permission = "coin_view"
	PermissionCoinCreate Permission = "coin_create"
	PermissionCoinUpdate Permission = "coin_update"

	PermissionStrategyView   Permission = "strategy_view"
	PermissionStrategyCreate Permission = "strategy_create"
	PermissionStrategyUpdate Permission = "strategy_update"

	PermissionTraderView   Permission = "trader_view"
	PermissionTraderCreate Permission = "trader_create"
	PermissionTraderUpdate Permission = "trader_update"
)

func PermissionTenantValues() []Permission {
	permissions := []Permission{
		PermissionSystemSetting,
		PermissionSystemAuditLog,

		PermissionRoleView,
		PermissionRoleCreate,
		PermissionRoleUpdate,

		PermissionUserView,
		PermissionUserCreate,
		PermissionUserUpdate,

		PermissionCoinView,
		PermissionCoinCreate,
		PermissionCoinUpdate,

		PermissionTraderView,
		PermissionTraderCreate,
		PermissionTraderUpdate,
	}
	return gopkg.UniqueFunc(slices.Sorted(slices.Values(permissions)), func(e Permission) Permission { return e })
}

func PermissionRootValues() []Permission {
	permissions := []Permission{
		PermissionSystemSetting,
		PermissionSystemAuditLog,

		PermissionClientView,
		PermissionClientCreate,
		PermissionClientDelete,

		PermissionRoleView,
		PermissionRoleCreate,
		PermissionRoleUpdate,

		PermissionUserView,
		PermissionUserCreate,
		PermissionUserUpdate,

		PermissionTenantView,
		PermissionTenantCreate,
		PermissionTenantUpdate,

		PermissionStrategyView,
		PermissionStrategyCreate,
		PermissionStrategyUpdate,

		PermissionCoinView,
		PermissionCoinCreate,
		PermissionCoinUpdate,

		PermissionTraderView,
		PermissionTraderCreate,
		PermissionTraderUpdate,
	}
	return gopkg.UniqueFunc(slices.Sorted(slices.Values(permissions)), func(e Permission) Permission { return e })
}

func PermissionValues() []Permission {
	permissions := slices.Concat(PermissionRootValues(), PermissionTenantValues())
	return gopkg.UniqueFunc(slices.Sorted(slices.Values(permissions)), func(e Permission) Permission { return e })
}

// DataStatus
type DataStatus string

const (
	DataStatusEnable  DataStatus = "enable"
	DataStatusDisable DataStatus = "disable"
)

func DataStatusValues() []DataStatus {
	return []DataStatus{DataStatusEnable, DataStatusDisable}
}

// DataAction
type DataAction string

const (
	DataActionCreate        DataAction = "create"
	DataActionUpdate        DataAction = "update"
	DataActionDelete        DataAction = "delete"
	DataActionResetPassword DataAction = "reset_password"
)

func DataActionValues() []DataAction {
	return []DataAction{DataActionCreate, DataActionUpdate, DataActionDelete, DataActionResetPassword}
}

// ServerStatus
type ServerStatus string

const (
	ServerStatusOn  ServerStatus = "on"
	ServerStatusOff ServerStatus = "off"
)

func ServerStatusValues() []ServerStatus {
	return []ServerStatus{ServerStatusOn, ServerStatusOff}
}

// Exchange
type Exchange string

const (
	ExchangeBinance Exchange = "binance"
	ExchangeBitget  Exchange = "bitget"
	ExchangeBybit   Exchange = "bybit"
	ExchangeCex     Exchange = "cex"
	ExchangeGate    Exchange = "gate"
	ExchangeMexc    Exchange = "mexc"
	ExchangeProbit  Exchange = "probit"
)

func ExchangeValues() []Exchange {
	return []Exchange{ExchangeBinance, ExchangeBitget, ExchangeBybit, ExchangeCex, ExchangeGate, ExchangeMexc, ExchangeProbit}
}

// MarketType
type MarketType string

const (
	MarketTypeSpot   MarketType = "spot"
	MarketTypeMargin MarketType = "margin"
)

func MarketTypeValues() []MarketType {
	return []MarketType{MarketTypeSpot, MarketTypeMargin}
}

// MarginMode
type MarginMode string

const (
	MarginModeCross    MarginMode = "cross"
	MarginModeIsolated MarginMode = "isolated"
)

func MarginModeValues() []MarginMode {
	return []MarginMode{MarginModeCross, MarginModeIsolated}
}

// OrderType
type OrderType string

const (
	OrderTypeLimit  OrderType = "limit"
	OrderTypeMarket OrderType = "market"
)

func OrderTypeValues() []OrderType {
	return []OrderType{OrderTypeLimit, OrderTypeMarket}
}

// OrderSide
type OrderSide string

const (
	OrderSideBuy  OrderSide = "buy"
	OrderSideSell OrderSide = "sell"
)

func OrderSideValues() []OrderSide {
	return []OrderSide{OrderSideBuy, OrderSideSell}
}

// OrderStatus
type OrderStatus string

const (
	OrderStatusOpen   OrderStatus = "open"
	OrderStatusClosed OrderStatus = "closed"
)

func OrderStatusValues() []OrderStatus {
	return []OrderStatus{OrderStatusOpen, OrderStatusClosed}
}
