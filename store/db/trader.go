package db

import (
	"app/pkg/encryption"
	"app/pkg/enum"
	"app/pkg/validate"
	"context"
	"time"

	"github.com/nhnghia272/gopkg"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TraderDomain struct {
	BaseDomain `json:"inline"`
	Name       *string          `json:"name,omitempty" validate:"omitempty"`
	Exchange   *enum.Exchange   `json:"exchange,omitempty" validate:"omitempty,exchange"`
	Sandbox    *bool            `json:"sandbox,omitempty" validate:"omitempty"`
	ApiKey     *string          `json:"api_key,omitempty" validate:"omitempty"`
	Secret     *string          `json:"secret,omitempty" validate:"omitempty"`
	Password   *string          `json:"password,omitempty" validate:"omitempty"`
	DataStatus *enum.DataStatus `json:"data_status,omitempty" validate:"omitempty,data_status"`
	TenantId   *enum.Tenant     `json:"tenant_id,omitempty" validate:"omitempty,len=24"`
}

func (s *TraderDomain) Validate() error {
	s.BeforeSave()
	return validate.New().Validate(s)
}

func (s TraderDomain) BaseDto() *TraderBaseDto {
	return &TraderBaseDto{
		ID:         SID(s.ID),
		Name:       gopkg.Value(s.Name),
		Exchange:   gopkg.Value(s.Exchange),
		Sandbox:    gopkg.Value(s.Sandbox),
		DataStatus: gopkg.Value(s.DataStatus),
	}
}

func (s TraderDomain) CmsDto() *TraderCmsDto {
	return &TraderCmsDto{
		ID:         SID(s.ID),
		Name:       gopkg.Value(s.Name),
		Exchange:   gopkg.Value(s.Exchange),
		Sandbox:    gopkg.Value(s.Sandbox),
		ApiKey:     encryption.Decrypt(gopkg.Value(s.ApiKey), string(gopkg.Value(s.TenantId))),
		Secret:     encryption.Decrypt(gopkg.Value(s.Secret), string(gopkg.Value(s.TenantId))),
		Password:   encryption.Decrypt(gopkg.Value(s.Password), string(gopkg.Value(s.TenantId))),
		DataStatus: gopkg.Value(s.DataStatus),
		UpdatedBy:  gopkg.Value(s.UpdatedBy),
		UpdatedAt:  gopkg.Value(s.UpdatedAt),
	}
}

func (s TraderDomain) ConfigDto() *TraderConfigDto {
	return &TraderConfigDto{
		TraderId: enum.Trader(SID(s.ID)),
		ApiKey:   encryption.Decrypt(gopkg.Value(s.ApiKey), string(gopkg.Value(s.TenantId))),
		Secret:   encryption.Decrypt(gopkg.Value(s.Secret), string(gopkg.Value(s.TenantId))),
		Password: encryption.Decrypt(gopkg.Value(s.Password), string(gopkg.Value(s.TenantId))),
	}
}

type TraderBaseDto struct {
	ID         string          `json:"trader_id" example:"671db9eca1f1b1bdbf3d4611"`
	Name       string          `json:"name" example:"Aloha"`
	Exchange   enum.Exchange   `json:"exchange"`
	Sandbox    bool            `json:"sandbox"`
	DataStatus enum.DataStatus `json:"data_status"`
}

type TraderCmsDto struct {
	ID         string          `json:"trader_id" example:"671db9eca1f1b1bdbf3d4611"`
	Name       string          `json:"name" example:"Aloha"`
	Exchange   enum.Exchange   `json:"exchange"`
	Sandbox    bool            `json:"sandbox"`
	ApiKey     string          `json:"api_key" example:"api_key"`
	Secret     string          `json:"secret" example:"secret"`
	Password   string          `json:"password" example:"password"`
	DataStatus enum.DataStatus `json:"data_status"`
	UpdatedBy  string          `json:"updated_by" example:"editor"`
	UpdatedAt  time.Time       `json:"updated_at" example:"2006-01-02T15:04:05Z"`
}

type TraderConfigDto struct {
	TraderId enum.Trader `json:"trader_id"`
	ApiKey   string      `json:"api_key"`
	Secret   string      `json:"secret"`
	Password string      `json:"password"`
}

type TraderCmsData struct {
	Name       string          `json:"name" validate:"required" example:"Aloha"`
	Exchange   enum.Exchange   `json:"exchange" validate:"required,exchange"`
	Sandbox    bool            `json:"sandbox" validate:"omitempty"`
	ApiKey     string          `json:"api_key" validate:"required"`
	Secret     string          `json:"secret" validate:"required"`
	Password   string          `json:"password" validate:"omitempty"`
	DataStatus enum.DataStatus `json:"data_status" validate:"required,data_status"`
}

func (s TraderCmsData) Domain(domain *TraderDomain) *TraderDomain {
	domain.Name = gopkg.Pointer(s.Name)
	domain.Exchange = gopkg.Pointer(s.Exchange)
	domain.Sandbox = gopkg.Pointer(s.Sandbox)
	domain.ApiKey = gopkg.Pointer(s.ApiKey)
	domain.Secret = gopkg.Pointer(s.Secret)
	domain.Password = gopkg.Pointer(s.Password)
	domain.DataStatus = gopkg.Pointer(s.DataStatus)
	return domain
}

type TraderQuery struct {
	Query
	Search     *string          `json:"search" form:"search" validate:"omitempty"`
	DataStatus *enum.DataStatus `json:"data_status" form:"data_status" validate:"omitempty,data_status"`
	TenantId   *enum.Tenant     `json:"tenant_id" form:"tenant_id" validate:"omitempty,len=24"`
}

func (s *TraderQuery) Build() *TraderQuery {
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

type TraderCmsQuery struct {
	Query
	Search     *string          `json:"search" form:"search" validate:"omitempty"`
	DataStatus *enum.DataStatus `json:"data_status" form:"data_status" validate:"omitempty,data_status"`
}

func (s *TraderCmsQuery) BuildCore() *TraderQuery {
	return &TraderQuery{Query: s.Query, Search: s.Search, DataStatus: s.DataStatus}
}

type trader struct {
	repo *repo
}

func newTrader(ctx context.Context, col *mongo.Collection) *trader {
	col.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "tenant_id", Value: 1}},
	})
	col.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "api_key", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	return &trader{newrepo(col)}
}

func (s trader) CollectionName() string { return s.repo.col.Name() }
func (s trader) Save(ctx context.Context, domain *TraderDomain, opts ...*options.UpdateOptions) (*TraderDomain, error) {
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

func (s trader) FindOne(ctx context.Context, filter M, opts ...*options.FindOneOptions) (*TraderDomain, error) {
	domain := &TraderDomain{}
	return domain, s.repo.FindOne(ctx, filter, &domain, opts...)
}

func (s trader) Count(ctx context.Context, q *TraderQuery, opts ...*options.CountOptions) int64 {
	return s.repo.CountDocuments(ctx, q.Build().Query, opts...)
}

func (s trader) FindAll(ctx context.Context, q *TraderQuery, opts ...*options.FindOptions) ([]*TraderDomain, error) {
	domains := make([]*TraderDomain, 0)
	return domains, s.repo.FindAll(ctx, q.Build().Query, &domains, opts...)
}

func (s trader) FindOneById(ctx context.Context, id string) (*TraderDomain, error) {
	return s.FindOne(ctx, M{"_id": OID(id)})
}
