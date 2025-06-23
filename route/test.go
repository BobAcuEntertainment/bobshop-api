package route

import (
	"app/pkg/ecode"
	"app/store/db"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

type test struct {
	*middleware
}

func init() {
	handlers = append(handlers, func(m *middleware, r *gin.Engine) {
		s := test{m}

		v1 := r.Group("/test/v1")
		v1.POST("", s.NoAuth(), s.v1_Test())
	})
}

// @Tags Test
// @Summary Test
// @Security BearerAuth
// @Param body body db.AuthLoginData true "body"
// @Success 200 {object} db.AuthLoginData
// @Router /test/v1 [post]
func (s test) v1_Test() gin.HandlerFunc {
	return func(c *gin.Context) {
		data := &db.AuthLoginData{}
		if err := c.ShouldBind(data); err != nil {
			c.Error(ecode.BadRequest.Desc(err))
			return
		}

		jv, _ := json.Marshal(data)
		s.ws.Broadcast(jv)

		c.JSON(http.StatusOK, data)
	}
}
