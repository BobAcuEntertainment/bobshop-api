package route

import (
	"app/pkg/coin"
	"app/pkg/ecode"
	"app/pkg/encryption"
	"app/pkg/enum"
	"app/store/db"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nhnghia272/gopkg"
)

type trader struct {
	*middleware
}

func init() {
	handlers = append(handlers, func(m *middleware, r *gin.Engine) {
		s := trader{m}

		v1cms := r.Group("/cms/v1/traders")
		v1cms.GET("", s.BearerAuth(enum.PermissionTraderView), s.v1cms_List())
		v1cms.POST("", s.BearerAuth(enum.PermissionTraderCreate), s.v1cms_Create())
		v1cms.PUT("/:trader_id", s.BearerAuth(enum.PermissionTraderUpdate), s.v1cms_Update())
	})
}

// @Tags Cms
// @Summary List Traders
// @Security BearerAuth
// @Param query query db.TraderCmsQuery false "query"
// @Success 200 {object} []db.TraderCmsDto
// @Failure 400,401,403,500 {object} ecode.Error
// @Router /cms/v1/traders [get]
func (s trader) v1cms_List() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := s.Session(c)

		qb := &db.TraderCmsQuery{}
		if err := c.ShouldBind(qb); err != nil {
			c.Error(ecode.BadRequest.Desc(err))
			return
		}

		query := qb.BuildCore()
		query.TenantId = gopkg.Pointer(session.TenantId)

		domains, _ := s.store.Db.Trader.FindAll(c.Request.Context(), query)
		results := gopkg.MapFunc(domains, func(domain *db.TraderDomain) *db.TraderCmsDto { return domain.CmsDto() })

		s.Pagination(c, s.store.Db.Trader.Count(c.Request.Context(), query), query.Page, query.Limit)

		c.JSON(http.StatusOK, results)
	}
}

// @Tags Cms
// @Summary Create Trader
// @Security BearerAuth
// @Param body body db.TraderCmsData true "body"
// @Success 200 {object} db.TraderCmsDto
// @Failure 400,401,403,409,500 {object} ecode.Error
// @Router /cms/v1/traders [post]
func (s trader) v1cms_Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := s.Session(c)

		data := &db.TraderCmsData{}
		if err := c.ShouldBind(data); err != nil {
			c.Error(ecode.BadRequest.Desc(err))
			return
		}

		data.ApiKey = encryption.Encrypt(data.ApiKey, string(session.TenantId))
		data.Secret = encryption.Encrypt(data.Secret, string(session.TenantId))
		data.Password = encryption.Encrypt(data.Password, string(session.TenantId))

		domain := data.Domain(&db.TraderDomain{})
		domain.TenantId = gopkg.Pointer(session.TenantId)
		domain.CreatedBy = gopkg.Pointer(session.Username)
		domain.UpdatedBy = gopkg.Pointer(session.Username)

		domain, err := s.store.Db.Trader.Save(c.Request.Context(), domain)
		if err != nil {
			c.Error(ecode.TraderConflict.Stack(err))
			return
		}

		s.AuditLog(c, s.store.Db.Trader.CollectionName(), enum.DataActionCreate, data, domain, db.SID(domain.ID))

		c.JSON(http.StatusOK, domain.CmsDto())
	}
}

// @Tags Cms
// @Summary Update Trader
// @Security BearerAuth
// @Param trader_id path string true "trader_id"
// @Param body body db.TraderCmsData true "body"
// @Success 200 {object} db.TraderCmsDto
// @Failure 400,401,403,404,409,500 {object} ecode.Error
// @Router /cms/v1/traders/{trader_id} [put]
func (s trader) v1cms_Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := s.Session(c)

		data := &db.TraderCmsData{}
		if err := c.ShouldBind(data); err != nil {
			c.Error(ecode.BadRequest.Desc(err))
			return
		}

		domain, err := s.store.Db.Trader.FindOneById(c.Request.Context(), c.Param("trader_id"))
		if err != nil || gopkg.Value(domain.TenantId) != session.TenantId {
			c.Error(ecode.TraderNotFound)
			return
		}

		data.ApiKey = encryption.Encrypt(data.ApiKey, string(session.TenantId))
		data.Secret = encryption.Encrypt(data.Secret, string(session.TenantId))
		data.Password = encryption.Encrypt(data.Password, string(session.TenantId))

		update := &db.TraderDomain{BaseDomain: db.BaseDomain{ID: domain.ID, UpdatedBy: gopkg.Pointer(session.Username)}}
		update = data.Domain(update)

		domain, err = s.store.Db.Trader.Save(c.Request.Context(), update)
		if err != nil {
			c.Error(ecode.TraderConflict.Stack(err))
			return
		}

		s.AuditLog(c, s.store.Db.Trader.CollectionName(), enum.DataActionUpdate, data, domain, db.SID(domain.ID))

		s.coin.DelCoin(&coin.TradeInfo{Exchange: gopkg.Value(domain.Exchange), Tenant: gopkg.Value(domain.TenantId), Trader: enum.Trader(db.SID(domain.ID)), Sandbox: gopkg.Value(domain.Sandbox)})

		if data.DataStatus == enum.DataStatusDisable {
			s.v1cms_UpdateConfig(c.Request.Context(), domain)
		}

		c.JSON(http.StatusOK, domain.CmsDto())
	}
}

func (s trader) v1cms_UpdateConfig(ctx context.Context, trader *db.TraderDomain) {
	domains, _ := s.store.Db.CoinConfig.FindAllByTraderId(ctx, db.SID(trader.ID))
	gopkg.LoopFunc(domains, func(domain *db.CoinConfigDomain) {
		domain, _ = s.store.Db.CoinConfig.Save(ctx, &db.CoinConfigDomain{BaseDomain: db.BaseDomain{ID: domain.ID}, DataStatus: gopkg.Pointer(enum.DataStatusDisable)})
		s.store.RunStrategy(ctx, domain)
	})
}
