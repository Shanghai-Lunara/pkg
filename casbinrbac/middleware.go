package casbinrbac

import (
	"github.com/Shanghai-Lunara/pkg/jwttoken"
	"github.com/Shanghai-Lunara/pkg/zaplogger"
	"github.com/gin-gonic/gin"
)

func Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get(jwttoken.TokenKey)
		if token == "" {
			c.Abort()
			return
		}
		tokenClaims, err := jwttoken.Parse(token)
		if err != nil {
			zaplogger.Sugar().Error(err)
			c.Abort()
			return
		}
		ok, err := Enforce(tokenClaims.Username, c.Param("namespace"), c.Param("permission"), c.Param("action"))
		if err != nil {
			zaplogger.Sugar().Error(err)
			c.Abort()
			return
		}
		if ok {
			c.Next()
		} else {
			c.Abort()
		}
	}
}
