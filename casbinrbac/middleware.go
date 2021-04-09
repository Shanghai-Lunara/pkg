package casbinrbac

import (
	"github.com/Shanghai-Lunara/pkg/zaplogger"
	"github.com/gin-gonic/gin"
)

func Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Token")
		_ = token
		user := "get from token"
		ok, err := Enforce(user, c.Param("namespace"), c.Param("permission"), c.Param("action"))
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
