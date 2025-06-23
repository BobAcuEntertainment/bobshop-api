package route

import (
	"app/env"
	"app/pkg/coin"
	"app/pkg/ecode"
	"app/pkg/enum"
	"app/pkg/mail"
	"app/pkg/oauth"
	"app/pkg/validate"
	"app/pkg/ws"
	"app/store"
	"app/store/db"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"slices"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-redis/redis_rate/v10"
	"github.com/nhnghia272/gopkg"
)

var handlers = make([]handler, 0)

type handler = func(*middleware, *gin.Engine)

func Bootstrap(store *store.Store) error {
	mdw := newMdw(store)
	gin.SetMode(gin.ReleaseMode)
	binding.Validator = validate.New()

	app := gin.New()
	app.NoRoute(mdw.NoRoute())
	app.Use(mdw.Cors(), mdw.Compress(), mdw.Logger(), mdw.Recover(), mdw.Error())

	for i := range handlers {
		handlers[i](mdw, app)
	}

	fmt.Println("Version: v1.0.142_update_logic")
	return app.Run(":" + env.Port)
}

type middleware struct {
	store   *store.Store
	oauth   *oauth.Oauth
	coin    *coin.Coins
	mail    *mail.Mail
	ws      *ws.Ws
	limiter *redis_rate.Limiter
}

func newMdw(store *store.Store) *middleware {
	return &middleware{store, oauth.New(store), coin.New(store), mail.New(store), ws.New(store), redis_rate.NewLimiter(store.Rdb.Instance())}
}

func (s middleware) NoRoute() gin.HandlerFunc {
	return func(c *gin.Context) { c.Error(ecode.ApiNotFound) }
}

func (s middleware) Cors() gin.HandlerFunc {
	config := cors.DefaultConfig()
	config.AllowWebSockets = true
	config.AllowAllOrigins = true
	config.AddAllowHeaders("Authorization", "Timezone")
	config.AddExposeHeaders("X-Pagination-Total", "X-Pagination-Page", "X-Pagination-Limit")
	return cors.New(config)
}

func (s middleware) Compress() gin.HandlerFunc {
	return gzip.Gzip(gzip.DefaultCompression)
}

func (s middleware) Logger() gin.HandlerFunc {
	return gin.Logger()
}

func (s middleware) Recover() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, err any) {
		s.ErrorFunc(c, err)
	})
}

func (s middleware) Error() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			s.ErrorFunc(c, c.Errors.Last().Err)
		}
	}
}

func (s middleware) ErrorFunc(c *gin.Context, err any) {
	switch e := err.(type) {
	case *ecode.Error:
		c.JSON(e.Status, e)
	default:
		err := ecode.InternalServerError.Stack(fmt.Errorf("%v", e))
		c.JSON(err.Status, err)
	}
}

func (s middleware) NoAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		session, err := s.oauth.NoAuth(c.Request)
		if err != nil {
			c.Error(err)
			c.Abort()
			return
		}
		if s.Limiter(c, nil) {
			c.Error(ecode.TooManyRequests)
			c.Abort()
			return
		}
		s.Session(c, session)
		c.Next()
	}
}

func (s middleware) BearerAuth(permissions ...enum.Permission) gin.HandlerFunc {
	return func(c *gin.Context) {
		session, err := s.oauth.BearerAuth(c.Request)
		if err != nil {
			c.Error(err)
			c.Abort()
			return
		}
		if len(permissions) > 0 && !slices.ContainsFunc(permissions, func(e enum.Permission) bool { return slices.Contains(session.Permissions, e) }) {
			c.Error(ecode.Forbidden)
			c.Abort()
			return
		}
		if s.Limiter(c, session) {
			c.Error(ecode.TooManyRequests)
			c.Abort()
			return
		}
		s.Session(c, session)
		c.Next()
	}
}

func (s middleware) BasicAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		session, err := s.oauth.BasicAuth(c.Request)
		if err != nil {
			c.Error(err)
			c.Abort()
			return
		}
		if s.Limiter(c, session) {
			c.Error(ecode.TooManyRequests)
			c.Abort()
			return
		}
		s.Session(c, session)
		c.Next()
	}
}

func (s middleware) Limiter(c *gin.Context, session *db.AuthSessionDto) bool {
	var (
		key   = c.Request.URL.Path + c.ClientIP()
		limit = redis_rate.PerMinute(100)
	)
	if session != nil {
		key = c.Request.URL.Path + string(session.TenantId)
		limit = redis_rate.Limit{Rate: 10_000, Burst: 1000, Period: time.Second * 10}
		if c.Request.Method != http.MethodGet {
			limit = redis_rate.Limit{Rate: 1_000, Burst: 100, Period: time.Second * 10}
		}
	}
	if res, err := s.limiter.Allow(c.Request.Context(), key, limit); err != nil || res.Allowed == 0 {
		return true
	}
	return false
}

func (s middleware) Session(c *gin.Context, session ...*db.AuthSessionDto) *db.AuthSessionDto {
	if len(session) == 0 {
		session, ok := c.Get(reflect.TypeOf(db.AuthSessionDto{}).Name())
		if !ok {
			return &db.AuthSessionDto{}
		}
		return session.(*db.AuthSessionDto)
	}
	c.Set(reflect.TypeOf(db.AuthSessionDto{}).Name(), session[0])
	return session[0]
}

func (s middleware) Pagination(c *gin.Context, total, page, limit int64) {
	c.Header("X-Pagination-Total", strconv.Itoa(int(total)))
	c.Header("X-Pagination-Page", strconv.Itoa(int(page)))
	c.Header("X-Pagination-Limit", strconv.Itoa(int(limit)))
}

func (s middleware) Permissions(session *db.AuthSessionDto) []enum.Permission {
	permissions := enum.PermissionRootValues()
	if session.IsTenant {
		permissions = enum.PermissionTenantValues()
	}
	return permissions
}

func (s middleware) AuditLog(c *gin.Context, name string, action enum.DataAction, data, domain any, domain_id string) {
	session := s.Session(c)
	byteData, _ := json.Marshal(data)
	byteDomain, _ := json.Marshal(domain)

	audit := &db.AuditLogDomain{}
	audit.Name = gopkg.Pointer(name)
	audit.Url = gopkg.Pointer(c.Request.URL.RequestURI())
	audit.Method = gopkg.Pointer(c.Request.Method)
	audit.Data = gopkg.Pointer(byteData)
	audit.Domain = gopkg.Pointer(byteDomain)
	audit.DomainId = gopkg.Pointer(domain_id)
	audit.Action = gopkg.Pointer(action)
	audit.TenantId = gopkg.Pointer(session.TenantId)
	audit.UpdatedBy = gopkg.Pointer(session.Username)

	s.store.Db.AuditLog.Save(c.Request.Context(), audit)
}
