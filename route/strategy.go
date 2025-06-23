package route

import (
	"app/pkg/ecode"
	"app/pkg/enum"
	"app/store/db"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nhnghia272/gopkg"
)

type strategy struct {
	*middleware
}

func init() {
	handlers = append(handlers, func(m *middleware, r *gin.Engine) {
		s := strategy{m}

		v1cms := r.Group("/cms/v1/strategies")
		v1cms.GET("", s.BearerAuth(enum.PermissionStrategyView), s.v1cms_List())
		v1cms.POST("", s.BearerAuth(enum.PermissionStrategyCreate), s.v1cms_Create())
		v1cms.PUT("/:strategy_id", s.BearerAuth(enum.PermissionStrategyUpdate), s.v1cms_Update())
	})
}

// @Tags Cms
// @Summary List Strategies
// @Security BearerAuth
// @Param query query db.StrategyCmsQuery false "query"
// @Success 200 {object} []db.StrategyCmsDto
// @Failure 400,401,403,500 {object} ecode.Error
// @Router /cms/v1/strategies [get]
func (s strategy) v1cms_List() gin.HandlerFunc {
	return func(c *gin.Context) {
		qb := &db.StrategyCmsQuery{}
		if err := c.ShouldBind(qb); err != nil {
			c.Error(ecode.BadRequest.Desc(err))
			return
		}

		query := qb.BuildCore()

		domains, _ := s.store.Db.Strategy.FindAll(c.Request.Context(), query)
		results := gopkg.MapFunc(domains, func(domain *db.StrategyDomain) *db.StrategyCmsDto { return domain.CmsDto() })

		s.Pagination(c, s.store.Db.Strategy.Count(c.Request.Context(), query), query.Page, query.Limit)

		c.JSON(http.StatusOK, results)
	}
}

// @Tags Cms
// @Summary Create Strategy
// @Security BearerAuth
// @Param body body db.StrategyCmsData true "body"
// @Success 200 {object} db.StrategyCmsDto
// @Failure 400,401,403,409,500 {object} ecode.Error
// @Router /cms/v1/strategies [post]
func (s strategy) v1cms_Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := s.Session(c)

		data := &db.StrategyCmsData{}
		if err := c.ShouldBind(data); err != nil {
			c.Error(ecode.BadRequest.Desc(err))
			return
		}

		domain := data.Domain(&db.StrategyDomain{})
		domain.CreatedBy = gopkg.Pointer(session.Username)
		domain.UpdatedBy = gopkg.Pointer(session.Username)

		domain, err := s.store.Db.Strategy.Save(c.Request.Context(), domain)
		if err != nil {
			c.Error(ecode.StrategyConflict.Stack(err))
			return
		}

		s.AuditLog(c, s.store.Db.Strategy.CollectionName(), enum.DataActionCreate, data, domain, db.SID(domain.ID))

		c.JSON(http.StatusOK, domain.CmsDto())
	}
}

// @Tags Cms
// @Summary Update Strategy
// @Security BearerAuth
// @Param strategy_id path string true "strategy_id"
// @Param body body db.StrategyCmsData true "body"
// @Success 200 {object} db.StrategyCmsDto
// @Failure 400,401,403,404,409,500 {object} ecode.Error
// @Router /cms/v1/strategies/{strategy_id} [put]
func (s strategy) v1cms_Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := s.Session(c)

		data := &db.StrategyCmsData{}
		if err := c.ShouldBind(data); err != nil {
			c.Error(ecode.BadRequest.Desc(err))
			return
		}

		domain, err := s.store.Db.Strategy.FindOneById(c.Request.Context(), c.Param("strategy_id"))
		if err != nil {
			c.Error(ecode.StrategyNotFound)
			return
		}

		update := &db.StrategyDomain{BaseDomain: db.BaseDomain{ID: domain.ID, UpdatedBy: gopkg.Pointer(session.Username)}}
		update = data.Domain(update)

		domain, err = s.store.Db.Strategy.Save(c.Request.Context(), update)
		if err != nil {
			c.Error(ecode.StrategyConflict.Stack(err))
			return
		}

		s.AuditLog(c, s.store.Db.Strategy.CollectionName(), enum.DataActionUpdate, data, domain, db.SID(domain.ID))

		if data.DataStatus == enum.DataStatusDisable {
			s.v1cms_UpdateConfig(c.Request.Context(), domain)
		}

		c.JSON(http.StatusOK, domain.CmsDto())
	}
}

func (s strategy) v1cms_UpdateConfig(ctx context.Context, strategy *db.StrategyDomain) {
	s.store.Db.CoinConfig.DisableByStrategyId(ctx, db.SID(strategy.ID))
	s.store.DelStrategy(ctx, db.SID(strategy.ID))
}
