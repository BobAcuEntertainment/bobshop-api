package db

import (
	"app/pkg/enum"
	"app/pkg/validate"
	"context"
	"encoding/json"
	"time"

	"github.com/nhnghia272/gopkg"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CoinConfigDomain struct {
	BaseDomain  `json:"inline"`
	Name        *string                    `json:"name,omitempty" validate:"omitempty"`
	Description *string                    `json:"description,omitempty" validate:"omitempty"`
	Exchange    *enum.Exchange             `json:"exchange,omitempty" validate:"omitempty,exchange"`
	Sandbox     *bool                      `json:"sandbox,omitempty" validate:"omitempty"`
	TraderIds   *[]string                  `json:"trader_ids,omitempty" validate:"omitempty,dive,len=24"`
	StrategyId  *string                    `json:"strategy_id,omitempty" validate:"omitempty,len=24"`
	MarketType  *enum.MarketType           `json:"market_type,omitempty" validate:"omitempty,market_type"`
	MarginMode  *enum.MarginMode           `json:"margin_mode,omitempty" validate:"omitempty,margin_mode"`
	Symbols     *[]string                  `json:"symbols,omitempty" validate:"omitempty,dive"`
	Attributes  *[]CoinConfigAttributeInfo `json:"attributes,omitempty" validate:"omitempty,dive"`
	DataStatus  *enum.DataStatus           `json:"data_status,omitempty" validate:"omitempty,data_status"`
	Server      *CoinConfigServerInfo      `json:"server,omitempty" validate:"omitempty"`
	TenantId    *enum.Tenant               `json:"tenant_id,omitempty" validate:"omitempty,len=24"`
}

func (s *CoinConfigDomain) Validate() error {
	s.BeforeSave()
	return validate.New().Validate(s)
}

func (s CoinConfigDomain) CmsDto() *CoinConfigCmsDto {
	return &CoinConfigCmsDto{
		ID:          SID(s.ID),
		Name:        gopkg.Value(s.Name),
		Description: gopkg.Value(s.Description),
		Exchange:    gopkg.Value(s.Exchange),
		Sandbox:     gopkg.Value(s.Sandbox),
		TraderIds:   gopkg.Value(s.TraderIds),
		StrategyId:  gopkg.Value(s.StrategyId),
		MarketType:  gopkg.Value(s.MarketType),
		MarginMode:  gopkg.Value(s.MarginMode),
		Symbols:     gopkg.Value(s.Symbols),
		Attributes:  gopkg.Value(s.Attributes),
		DataStatus:  gopkg.Value(s.DataStatus),
		Server:      gopkg.Value(s.Server),
		UpdatedBy:   gopkg.Value(s.UpdatedBy),
		UpdatedAt:   gopkg.Value(s.UpdatedAt),
	}
}

func (s CoinConfigDomain) StrategyDto() *CoinConfigStrategyDto {
	dto := &CoinConfigStrategyDto{
		ConfigId:    SID(s.ID),
		Exchange:    gopkg.Value(s.Exchange),
		TenantId:    gopkg.Value(s.TenantId),
		StrategyId:  enum.Strategy(gopkg.Value(s.StrategyId)),
		Status:      gopkg.Value(s.DataStatus),
		Sandbox:     gopkg.Value(s.Sandbox),
		Attributes:  make(map[string]any),
		Credentials: make([]CoinConfigCredentialInfo, 0),
	}

	gopkg.LoopFunc(gopkg.Value(s.Attributes), func(attribute CoinConfigAttributeInfo) {
		dto.Attributes[attribute.Key] = attribute.Value
	})

	dto.Attributes["market_type"] = gopkg.Value(s.MarketType)
	dto.Attributes["margin_mode"] = gopkg.Value(s.MarginMode)
	dto.Attributes["symbols"] = gopkg.Value(s.Symbols)

	return dto
}

type CoinConfigCmsDto struct {
	ID          string                    `json:"config_id" example:"671dfc49f06ba89b1821cc5a"`
	Name        string                    `json:"name"`
	Description string                    `json:"description"`
	Exchange    enum.Exchange             `json:"exchange" `
	Sandbox     bool                      `json:"sandbox"`
	TraderIds   []string                  `json:"trader_ids" example:"671db9eca1f1b1bdbf3d4611"`
	StrategyId  string                    `json:"strategy_id" example:"671db9eca1f1b1bdbf3d4627"`
	MarketType  enum.MarketType           `json:"market_type"`
	MarginMode  enum.MarginMode           `json:"margin_mode"`
	Symbols     []string                  `json:"symbols" example:"BTC/USDT"`
	Attributes  []CoinConfigAttributeInfo `json:"attributes"`
	DataStatus  enum.DataStatus           `json:"data_status"`
	Server      CoinConfigServerInfo      `json:"server"`
	UpdatedBy   string                    `json:"updated_by" example:"editor"`
	UpdatedAt   time.Time                 `json:"updated_at" example:"2006-01-02T15:04:05Z"`
}

type CoinConfigStrategyDto struct {
	ConfigId    string                     `json:"config_id"`
	Exchange    enum.Exchange              `json:"exchange"`
	TenantId    enum.Tenant                `json:"tenant_id"`
	StrategyId  enum.Strategy              `json:"strategy_id"`
	Status      enum.DataStatus            `json:"status"`
	Sandbox     bool                       `json:"sandbox"`
	Attributes  map[string]any             `json:"attributes"`
	Credentials []CoinConfigCredentialInfo `json:"credentials"`
}

func (s CoinConfigStrategyDto) Bytes() []byte {
	bytes, _ := json.Marshal(s)
	return bytes
}

type CoinConfigCmsData struct {
	Name        string                    `json:"name" validate:"required"`
	Description string                    `json:"description" validate:"required"`
	Exchange    enum.Exchange             `json:"exchange" validate:"required,exchange"`
	Sandbox     bool                      `json:"sandbox" validate:"omitempty"`
	TraderIds   []string                  `json:"trader_ids" validate:"required,min=1,dive,len=24" example:"671db9eca1f1b1bdbf3d4611"`
	StrategyId  string                    `json:"strategy_id" validate:"required,len=24" example:"671db9eca1f1b1bdbf3d4627"`
	MarketType  enum.MarketType           `json:"market_type" validate:"required,market_type"`
	MarginMode  enum.MarginMode           `json:"margin_mode" validate:"required,margin_mode"`
	Symbols     []string                  `json:"symbols" validate:"required,min=1,dive" example:"BTC/USDT"`
	Attributes  []CoinConfigAttributeInfo `json:"attributes" validate:"required,dive"`
	DataStatus  enum.DataStatus           `json:"data_status" validate:"required,data_status"`
}

func (s CoinConfigCmsData) Domain(domain *CoinConfigDomain) *CoinConfigDomain {
	domain.Name = gopkg.Pointer(s.Name)
	domain.Description = gopkg.Pointer(s.Description)
	domain.Exchange = gopkg.Pointer(s.Exchange)
	domain.Sandbox = gopkg.Pointer(s.Sandbox)
	domain.TraderIds = gopkg.Pointer(s.TraderIds)
	domain.StrategyId = gopkg.Pointer(s.StrategyId)
	domain.MarketType = gopkg.Pointer(s.MarketType)
	domain.MarginMode = gopkg.Pointer(s.MarginMode)
	domain.Symbols = gopkg.Pointer(s.Symbols)
	domain.Attributes = gopkg.Pointer(s.Attributes)
	domain.DataStatus = gopkg.Pointer(s.DataStatus)
	return domain
}

type CoinConfigAttributeInfo struct {
	Key   string `json:"key" validate:"required"`
	Value any    `json:"value" validate:"required"`
}

type CoinConfigCredentialInfo struct {
	TraderId enum.Trader `json:"trader_id"`
	ApiKey   string      `json:"api_key"`
	Secret   string      `json:"secret"`
	Password string      `json:"password"`
}

type CoinConfigServerInfo struct {
	Message string            `json:"message" validate:"required"`
	Status  enum.ServerStatus `json:"status" validate:"required,server_status"`
}

type CoinConfigQuery struct {
	Query
	Search     *string          `json:"search" form:"search" validate:"omitempty"`
	DataStatus *enum.DataStatus `json:"data_status" form:"data_status" validate:"omitempty,data_status"`
	TenantId   *enum.Tenant     `json:"tenant_id" form:"tenant_id" validate:"omitempty,len=24"`
}

func (s *CoinConfigQuery) Build() *CoinConfigQuery {
	if s.Filter == nil {
		s.Filter = M{}
	}
	if s.Search != nil {
		s.Filter["$or"] = []M{{"name": Regex(gopkg.Value(s.Search))}}
	}
	if s.DataStatus != nil {
		s.Filter["data_status"] = s.DataStatus
	}
	if s.TenantId != nil {
		s.Filter["tenant_id"] = s.TenantId
	}
	return s
}

type CoinConfigCmsQuery struct {
	Query
	Search     *string          `json:"search" form:"search" validate:"omitempty"`
	DataStatus *enum.DataStatus `json:"data_status" form:"data_status" validate:"omitempty,data_status"`
}

func (s *CoinConfigCmsQuery) BuildCore() *CoinConfigQuery {
	return &CoinConfigQuery{Query: s.Query, Search: s.Search, DataStatus: s.DataStatus}
}

type coin_config struct {
	repo *repo
}

func newCoinConfig(ctx context.Context, col *mongo.Collection) *coin_config {
	col.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "tenant_id", Value: 1}},
	})
	return &coin_config{newrepo(col)}
}

func (s coin_config) CollectionName() string { return s.repo.col.Name() }

func (s coin_config) Save(ctx context.Context, domain *CoinConfigDomain, opts ...*options.UpdateOptions) (*CoinConfigDomain, error) {
	if err := domain.Validate(); err != nil {
		return nil, err
	}
	id, err := s.repo.Save(ctx, domain.ID, domain, opts...)
	if err != nil {
		return nil, err
	}
	domain.ID = id
	return s.FindOneById(ctx, SID(id))
}

func (s coin_config) UpdateOne(ctx context.Context, filter M, update M, opts ...*options.UpdateOptions) error {
	return s.repo.UpdateOne(ctx, filter, update, opts...)
}

func (s coin_config) FindOne(ctx context.Context, filter M, opts ...*options.FindOneOptions) (*CoinConfigDomain, error) {
	domain := &CoinConfigDomain{}
	return domain, s.repo.FindOne(ctx, filter, &domain, opts...)
}

func (s coin_config) Count(ctx context.Context, q *CoinConfigQuery, opts ...*options.CountOptions) int64 {
	return s.repo.CountDocuments(ctx, q.Build().Query, opts...)
}

func (s coin_config) FindAll(ctx context.Context, q *CoinConfigQuery, opts ...*options.FindOptions) ([]*CoinConfigDomain, error) {
	domains := make([]*CoinConfigDomain, 0)
	return domains, s.repo.FindAll(ctx, q.Build().Query, &domains, opts...)
}

func (s coin_config) FindOneById(ctx context.Context, id string) (*CoinConfigDomain, error) {
	return s.FindOne(ctx, M{"_id": OID(id)})
}

func (s coin_config) FindAllByTraderId(ctx context.Context, traderId string) ([]*CoinConfigDomain, error) {
	return s.FindAll(ctx, &CoinConfigQuery{Query: Query{Filter: M{"trader_ids": traderId}}})
}

func (s coin_config) FindAllByStrategyId(ctx context.Context, strategyId string) ([]*CoinConfigDomain, error) {
	return s.FindAll(ctx, &CoinConfigQuery{Query: Query{Filter: M{"strategy_id": strategyId}}})
}

func (s coin_config) DisableByTraderId(ctx context.Context, traderId string) error {
	return s.UpdateOne(ctx, M{"trader_ids": traderId}, M{"$set": M{"data_status": enum.DataStatusDisable}})
}

func (s coin_config) DisableByStrategyId(ctx context.Context, strategyId string) error {
	return s.UpdateOne(ctx, M{"strategy_id": strategyId}, M{"$set": M{"data_status": enum.DataStatusDisable}})
}
