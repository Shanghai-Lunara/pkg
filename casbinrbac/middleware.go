package casbinrbac

import (
	"fmt"
	"github.com/Shanghai-Lunara/pkg/jwttoken"
	"github.com/Shanghai-Lunara/pkg/zaplogger"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get(jwttoken.TokenKey)
		if token == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		tokenClaims, err := jwttoken.Parse(token)
		if err != nil {
			zaplogger.Sugar().Error(err)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		ok, err := Enforce(tokenClaims.Username, c.Param("namespace"), c.Param("permission"), c.Param("action"))
		if err != nil {
			zaplogger.Sugar().Error(err)
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		if ok {
			c.Next()
		} else {
			c.AbortWithStatus(http.StatusForbidden)
		}
	}
}

func FullPath(relativePath, suffixPath string) string {
	if relativePath == "" || relativePath == "/" {
		return suffixPath
	} else {
		return fmt.Sprintf("/%s%s", relativePath, suffixPath)
	}
}
