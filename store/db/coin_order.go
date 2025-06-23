package db

import (
	"app/pkg/enum"
	"app/pkg/tz"
	"app/pkg/validate"
	"context"
	"time"

	"github.com/nhnghia272/gopkg"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CoinOrderDomain struct {
	BaseDomain `json:"inline"`
	ConfigId   *string         `json:"config_id,omitempty" validate:"omitempty,len=24"`
	TraderId   *string         `json:"trader_id,omitempty" validate:"omitempty,len=24"`
	Order      *map[string]any `json:"order,omitempty" validate:"omitempty"`
	TenantId   *enum.Tenant    `json:"tenant_id,omitempty" validate:"omitempty,len=24"`
}

func (s *CoinOrderDomain) Validate() error {
	s.BeforeSave()
	return validate.New().Validate(s)
}

func (s CoinOrderDomain) CmsDto() *CoinOrderCmsDto {
	return &CoinOrderCmsDto{
		ID:        SID(s.ID),
		ConfigId:  gopkg.Value(s.ConfigId),
		TraderId:  gopkg.Value(s.TraderId),
		Order:     gopkg.Value(s.Order),
		UpdatedBy: gopkg.Value(s.UpdatedBy),
		UpdatedAt: gopkg.Value(s.UpdatedAt),
	}
}

type CoinOrderCmsDto struct {
	ID        string         `json:"order_id" example:"671dfc49f06ba89b1821cc5a"`
	ConfigId  string         `json:"config_id" example:"671dfc49f06ba89b1821cc5a"`
	TraderId  string         `json:"trader_id" example:"671dfc49f06ba89b1821cc5a"`
	Order     map[string]any `json:"order"`
	UpdatedBy string         `json:"updated_by" example:"editor"`
	UpdatedAt time.Time      `json:"updated_at" example:"2006-01-02T15:04:05Z"`
}

type CoinOrderCmsReportDto struct {
	Summary24H CoinOrderCmsReportSummary24HInfo `json:"summary_24h"`
	Summary    CoinOrderCmsReportSummaryInfo    `json:"summary"`
	Charts     []CoinOrderCmsReportChartInfo    `json:"charts"`
}

type CoinOrderCmsReportSummary24HInfo struct {
	TotalFee            float64 `json:"total_fee"`
	TotalOrder          float64 `json:"total_order"`
	TotalToken          float64 `json:"total_token"`
	TotalVolume         float64 `json:"total_volume"`
	TotalTokenExchange  float64 `json:"total_token_exchange"`
	TotalVolumeExchange float64 `json:"total_volume_exchange"`
}

type CoinOrderCmsReportSummaryInfo struct {
	TotalFee    float64 `json:"total_fee"`
	TotalOrder  float64 `json:"total_order"`
	TotalToken  float64 `json:"total_token"`
	TotalVolume float64 `json:"total_volume"`
}

type CoinOrderCmsReportChartInfo struct {
	Name        string  `json:"name"`
	TotalFee    float64 `json:"total_fee"`
	TotalOrder  float64 `json:"total_order"`
	TotalToken  float64 `json:"total_token"`
	TotalVolume float64 `json:"total_volume"`
}

type CoinOrderAlertData struct {
	ConfigId string         `json:"config_id" validate:"required,len=24"`
	TraderId string         `json:"trader_id" validate:"required,len=24"`
	TenantId enum.Tenant    `json:"tenant_id" validate:"required,len=24"`
	Order    map[string]any `json:"order" validate:"required"`
}

func (s CoinOrderAlertData) Domain(domain *CoinOrderDomain) *CoinOrderDomain {
	domain.ConfigId = gopkg.Pointer(s.ConfigId)
	domain.TraderId = gopkg.Pointer(s.TraderId)
	domain.TenantId = gopkg.Pointer(s.TenantId)
	domain.Order = gopkg.Pointer(s.Order)
	return domain
}

type CoinOrderQuery struct {
	Query
	Search    *string           `json:"search" form:"search" validate:"omitempty"`
	TraderId  *string           `json:"trader_id" form:"trader_id" validate:"omitempty,len=24"`
	ConfigId  *string           `json:"config_id" form:"config_id" validate:"omitempty,len=24"`
	Status    *enum.OrderStatus `json:"status" form:"status" validate:"omitempty,order_status"`
	StartDate *time.Time        `json:"start_date" form:"start_date" validate:"omitempty"`
	EndDate   *time.Time        `json:"end_date" form:"end_date" validate:"omitempty"`
	Timezone  *tz.Timezone      `json:"timezone" form:"timezone" validate:"omitempty"`
	TenantId  *enum.Tenant      `json:"tenant_id" form:"tenant_id" validate:"omitempty,len=24"`
}

func (s *CoinOrderQuery) Build() *CoinOrderQuery {
	if s.Filter == nil {
		s.Filter = M{}
	}
	if s.Search != nil {
		s.Filter["$or"] = []M{{"name": Regex(gopkg.Value(s.Search))}}
	}
	if s.TraderId != nil {
		s.Filter["trader_id"] = s.TraderId
	}
	if s.ConfigId != nil {
		s.Filter["config_id"] = s.ConfigId
	}
	if s.Status != nil {
		s.Filter["order.status"] = s.Status
	}
	if s.StartDate != nil && s.EndDate != nil {
		s.Filter["created_at"] = M{"$gte": s.StartDate, "$lte": s.EndDate}
	}
	if s.TenantId != nil {
		s.Filter["tenant_id"] = s.TenantId
	}
	return s
}

type CoinOrderCmsQuery struct {
	Query
	Search    *string           `json:"search" form:"search" validate:"omitempty"`
	TraderId  *string           `json:"trader_id" form:"trader_id" validate:"omitempty,len=24"`
	ConfigId  *string           `json:"config_id" form:"config_id" validate:"omitempty,len=24"`
	Status    *enum.OrderStatus `json:"status" form:"status" validate:"omitempty,order_status"`
	StartDate *string           `json:"start_date" form:"start_date" validate:"required,datetime=2006-01-02" example:"2006-01-02"`
	EndDate   *string           `json:"end_date" form:"end_date" validate:"required,datetime=2006-01-02" example:"2006-01-02"`
}

func (s *CoinOrderCmsQuery) BuildCore(timezone tz.Timezone) *CoinOrderQuery {
	var (
		end   = tz.EndOfDay(tz.ParseDate(gopkg.Value(s.EndDate), timezone))
		start = tz.StartOfDay(tz.ParseDate(gopkg.Value(s.StartDate), timezone))
	)
	return &CoinOrderQuery{Query: s.Query, Search: s.Search, TraderId: s.TraderId, ConfigId: s.ConfigId, Status: s.Status, StartDate: gopkg.Pointer(start), EndDate: gopkg.Pointer(end), Timezone: gopkg.Pointer(timezone)}
}

type CoinOrderCmsReportQuery struct {
	TraderId  *string `json:"trader_id" form:"trader_id" validate:"omitempty,len=24"`
	ConfigId  *string `json:"config_id" form:"config_id" validate:"omitempty,len=24"`
	StartDate *string `json:"start_date" form:"start_date" validate:"required,datetime=2006-01-02" example:"2006-01-02"`
	EndDate   *string `json:"end_date" form:"end_date" validate:"required,datetime=2006-01-02" example:"2006-01-02"`
}

func (s *CoinOrderCmsReportQuery) BuildCore(timezone tz.Timezone) *CoinOrderQuery {
	var (
		end   = tz.EndOfDay(tz.ParseDate(gopkg.Value(s.EndDate), timezone))
		start = tz.StartOfDay(tz.ParseDate(gopkg.Value(s.StartDate), timezone))
	)
	return &CoinOrderQuery{TraderId: s.TraderId, ConfigId: s.ConfigId, StartDate: gopkg.Pointer(start), EndDate: gopkg.Pointer(end), Timezone: gopkg.Pointer(timezone)}
}

type coin_order struct {
	repo *repo
}

func newCoinOrder(ctx context.Context, col *mongo.Collection) *coin_order {
	col.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "tenant_id", Value: 1}},
	})
	col.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "config_id", Value: 1}, {Key: "order.id", Value: 1}, {Key: "order.status", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	return &coin_order{newrepo(col)}
}

func (s coin_order) CollectionName() string { return s.repo.col.Name() }

func (s coin_order) Save(ctx context.Context, domain *CoinOrderDomain, opts ...*options.UpdateOptions) (*CoinOrderDomain, error) {
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

func (s coin_order) UpdateOne(ctx context.Context, filter M, update M, opts ...*options.UpdateOptions) error {
	return s.repo.UpdateOne(ctx, filter, update, opts...)
}

func (s coin_order) FindOne(ctx context.Context, filter M, opts ...*options.FindOneOptions) (*CoinOrderDomain, error) {
	domain := &CoinOrderDomain{}
	return domain, s.repo.FindOne(ctx, filter, &domain, opts...)
}

func (s coin_order) Count(ctx context.Context, q *CoinOrderQuery, opts ...*options.CountOptions) int64 {
	return s.repo.CountDocuments(ctx, q.Build().Query, opts...)
}

func (s coin_order) FindAll(ctx context.Context, q *CoinOrderQuery, opts ...*options.FindOptions) ([]*CoinOrderDomain, error) {
	domains := make([]*CoinOrderDomain, 0)
	return domains, s.repo.FindAll(ctx, q.Build().Query, &domains, opts...)
}

func (s coin_order) FindOneById(ctx context.Context, id string) (*CoinOrderDomain, error) {
	return s.FindOne(ctx, M{"_id": OID(id)})
}

func (s coin_order) ReportCoinOrderDetails(ctx context.Context, q *CoinOrderQuery) (*CoinOrderCmsReportDto, error) {
	var (
		summary    = CoinOrderCmsReportSummaryInfo{}
		summary24H = CoinOrderCmsReportSummary24HInfo{}
		charts     = make([]CoinOrderCmsReportChartInfo, 0)
		charts24h  = make([]CoinOrderCmsReportChartInfo, 0)
		format     = "%H"
		match      = q.Build().Filter
	)

	if q.StartDate != nil && q.EndDate != nil && gopkg.Value(q.EndDate).Sub(gopkg.Value(q.StartDate)) > 24*time.Hour {
		format = "%Y-%m-%d"
	}

	pipeline := []M{
		{"$match": match},
		{"$group": M{"_id": M{"$dateToString": M{"format": format, "date": "$created_at", "timezone": q.Timezone}},
			"total_fee":    M{"$sum": "$order.fee_convert"},
			"total_order":  M{"$sum": 1},
			"total_token":  M{"$sum": "$order.amount"},
			"total_volume": M{"$sum": "$order.cost"},
		}},
		{"$sort": M{"_id": 1}},
		{"$project": M{"name": "$_id", "total_fee": 1, "total_order": 1, "total_token": 1, "total_volume": 1}},
	}

	if err := s.repo.Aggregate(ctx, pipeline, &charts, options.Aggregate().SetAllowDiskUse(true)); err != nil {
		return nil, err
	}

	match["created_at"] = M{"$gte": time.Now().In(tz.Location(gopkg.Value(q.Timezone))).Add(-24 * time.Hour)}

	pipeline24h := []M{
		{"$match": match},
		{"$group": M{"_id": nil,
			"total_fee":    M{"$sum": "$order.fee_convert"},
			"total_order":  M{"$sum": 1},
			"total_token":  M{"$sum": "$order.amount"},
			"total_volume": M{"$sum": "$order.cost"},
		}},
	}

	if err := s.repo.Aggregate(ctx, pipeline24h, &charts24h, options.Aggregate().SetAllowDiskUse(true)); err != nil {
		return nil, err
	}

	if len(charts24h) > 0 {
		summary24H.TotalFee = charts24h[0].TotalFee
		summary24H.TotalOrder = charts24h[0].TotalOrder
		summary24H.TotalToken = charts24h[0].TotalToken
		summary24H.TotalVolume = charts24h[0].TotalVolume
	}

	for _, chart := range charts {
		summary.TotalFee += chart.TotalFee
		summary.TotalOrder += chart.TotalOrder
		summary.TotalToken += chart.TotalToken
		summary.TotalVolume += chart.TotalVolume
	}

	return &CoinOrderCmsReportDto{Summary24H: summary24H, Summary: summary, Charts: charts}, nil
}
