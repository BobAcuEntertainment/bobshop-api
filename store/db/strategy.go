package db

import (
	"app/pkg/enum"
	"app/pkg/validate"
	"context"
	"time"

	"github.com/nhnghia272/gopkg"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type StrategyDomain struct {
	BaseDomain  `json:"inline"`
	Name        *string                  `json:"name,omitempty" validate:"omitempty"`
	Description *string                  `json:"description,omitempty" validate:"omitempty"`
	Filename    *string                  `json:"filename,omitempty" validate:"omitempty,lowercase"`
	Attributes  *[]StrategyAttributeInfo `json:"attributes,omitempty" validate:"omitempty,dive"`
	DataStatus  *enum.DataStatus         `json:"data_status,omitempty" validate:"omitempty,data_status"`
}

func (s *StrategyDomain) Validate() error {
	s.BeforeSave()
	return validate.New().Validate(s)
}

func (s StrategyDomain) CmsDto() *StrategyCmsDto {
	return &StrategyCmsDto{
		ID:          SID(s.ID),
		Name:        gopkg.Value(s.Name),
		Description: gopkg.Value(s.Description),
		Filename:    gopkg.Value(s.Filename),
		Attributes:  gopkg.Value(s.Attributes),
		DataStatus:  gopkg.Value(s.DataStatus),
		UpdatedBy:   gopkg.Value(s.UpdatedBy),
		UpdatedAt:   gopkg.Value(s.UpdatedAt),
	}
}

type StrategyCmsDto struct {
	ID          string                  `json:"strategy_id" example:"671db9eca1f1b1bdbf3d4627"`
	Name        string                  `json:"name" example:"Aloha"`
	Description string                  `json:"description" example:"Description"`
	Filename    string                  `json:"filename" example:"strategy_1.py"`
	Attributes  []StrategyAttributeInfo `json:"attributes"`
	DataStatus  enum.DataStatus         `json:"data_status"`
	UpdatedBy   string                  `json:"updated_by" example:"editor"`
	UpdatedAt   time.Time               `json:"updated_at" example:"2006-01-02T15:04:05Z"`
}

type StrategyCmsData struct {
	Name        string                  `json:"name" validate:"required" example:"Aloha"`
	Description string                  `json:"description" validate:"required" example:"Description"`
	Filename    string                  `json:"filename" validate:"required" example:"strategy_1.py"`
	Attributes  []StrategyAttributeInfo `json:"attributes" validate:"required,dive"`
	DataStatus  enum.DataStatus         `json:"data_status" validate:"required,data_status"`
}

func (s StrategyCmsData) Domain(domain *StrategyDomain) *StrategyDomain {
	domain.Name = gopkg.Pointer(s.Name)
	domain.Description = gopkg.Pointer(s.Description)
	domain.Filename = gopkg.Pointer(s.Filename)
	domain.Attributes = gopkg.Pointer(s.Attributes)
	domain.DataStatus = gopkg.Pointer(s.DataStatus)
	return domain
}

type StrategyAttributeInfo struct {
	Key     string `json:"key" validate:"required"`
	Type    string `json:"type" validate:"required"`
	Title   string `json:"title" validate:"required"`
	Default any    `json:"default" validate:"omitempty"`
	Values  any    `json:"values" validate:"omitempty"`
}

type StrategyQuery struct {
	Query
	Search     *string          `json:"search" form:"search" validate:"omitempty"`
	DataStatus *enum.DataStatus `json:"data_status" form:"data_status" validate:"omitempty,data_status"`
}

func (s *StrategyQuery) Build() *StrategyQuery {
	if s.Filter == nil {
		s.Filter = M{}
	}
	if s.Search != nil {
		s.Filter["$or"] = []M{{"name": Regex(gopkg.Value(s.Search))}}
	}
	if s.DataStatus != nil {
		s.Filter["data_status"] = s.DataStatus
	}
	return s
}

type StrategyCmsQuery struct {
	Query
	Search     *string          `json:"search" form:"search" validate:"omitempty"`
	DataStatus *enum.DataStatus `json:"data_status" form:"data_status" validate:"omitempty,data_status"`
}

func (s *StrategyCmsQuery) BuildCore() *StrategyQuery {
	return &StrategyQuery{Query: s.Query, Search: s.Search, DataStatus: s.DataStatus}
}

type strategy struct {
	repo *repo
}

func newStrategy(ctx context.Context, col *mongo.Collection) *strategy {
	col.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "filename", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	return &strategy{newrepo(col)}
}

func (s strategy) CollectionName() string { return s.repo.col.Name() }

func (s strategy) Save(ctx context.Context, domain *StrategyDomain, opts ...*options.UpdateOptions) (*StrategyDomain, error) {
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

func (s strategy) FindOne(ctx context.Context, filter M, opts ...*options.FindOneOptions) (*StrategyDomain, error) {
	domain := &StrategyDomain{}
	return domain, s.repo.FindOne(ctx, filter, &domain, opts...)
}

func (s strategy) Count(ctx context.Context, q *StrategyQuery, opts ...*options.CountOptions) int64 {
	return s.repo.CountDocuments(ctx, q.Build().Query, opts...)
}

func (s strategy) FindAll(ctx context.Context, q *StrategyQuery, opts ...*options.FindOptions) ([]*StrategyDomain, error) {
	domains := make([]*StrategyDomain, 0)
	return domains, s.repo.FindAll(ctx, q.Build().Query, &domains, opts...)
}

func (s strategy) FindOneById(ctx context.Context, id string) (*StrategyDomain, error) {
	return s.FindOne(ctx, M{"_id": OID(id)})
}

func (s strategy) FindOneByFilename(ctx context.Context, filename string) (*StrategyDomain, error) {
	return s.FindOne(ctx, M{"filename": filename})
}
