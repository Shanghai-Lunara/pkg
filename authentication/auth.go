package authentication

import (
	"github.com/Shanghai-Lunara/pkg/casbinrbac"
	"github.com/Shanghai-Lunara/pkg/jwttoken"
	"github.com/Shanghai-Lunara/pkg/zaplogger"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Authentication struct {
	relativePath string
	mysql        *casbinrbac.MysqlClusterPool
}

func New(relativePath string, router *gin.Engine) *Authentication {
	authentication = &Authentication{
		relativePath: relativePath,
		mysql:        casbinrbac.GetMysqlCluster(),
	}
	register(router.Group(relativePath))
	return authentication
}

var authentication *Authentication

const (
	AuthAccountLogin   = "/account/login"
	AuthAccountList    = "/account/list"
	AuthAccountAdd     = "/account/add"
	AuthAccountReset   = "/account/reset"
	AuthAccountDisable = "/account/disable/:account"
	AuthAccountEnable  = "/account/enable/:account"
)

const (
	ParamAccount  = "account"
	ParamPassword = "pwd"
)

func register(router *gin.RouterGroup) {
	router.Use(authentication.middleware())
	router.POST(AuthAccountLogin, authentication.login)
	router.GET(AuthAccountList, authentication.list)
	router.POST(AuthAccountAdd, authentication.add)
	router.POST(AuthAccountReset, authentication.reset)
	router.GET(AuthAccountDisable, authentication.disable)
	router.GET(AuthAccountEnable, authentication.enable)
}

func (a *Authentication) middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		switch c.FullPath() {
		case casbinrbac.FullPath(a.relativePath, AuthAccountLogin):
			c.Next()
		default:
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
			switch tokenClaims.IsAdmin {
			case true:
				c.Next()
			default:
				c.AbortWithStatus(http.StatusForbidden)
			}
		}
	}
}

func IsAdmin(username string) bool {
	switch username {
	case "admin":
		return true
	default:
		return false
	}
}

func (a *Authentication) login(c *gin.Context) {
	req := &LoginRequest{}
	if c.ShouldBindJSON(req) != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	}
	acc, err := Query(a.mysql.Slave, req.Account, req.Password)
	if err != nil {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	token, err := jwttoken.Generate(acc.Account, IsAdmin(acc.Account))
	if err != nil {
		zaplogger.Sugar().Errorw("jwttoken.Generate", "req", req, "account", acc, "err", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, &LoginResponse{Token: token, IsAdmin: IsAdmin(acc.Account)})
}

func (a *Authentication) list(c *gin.Context) {
	accountList, err := List(a.mysql.Slave)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, &ListResponse{
		Items: accountList,
	})
}

func (a *Authentication) add(c *gin.Context) {
	req := &AccountRequest{}
	if c.ShouldBindJSON(req) != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if Add(a.mysql.Master, req.Account, req.Password) != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, &BoolResultResponse{
		Result: true,
	})
}

func (a *Authentication) reset(c *gin.Context) {
	req := &AccountRequest{}
	if c.ShouldBindJSON(req) != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if ResetPassword(a.mysql.Master, req.Account, req.Password) != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, &BoolResultResponse{Result: true})
}

func (a *Authentication) disable(c *gin.Context) {
	if Operator(a.mysql.Master, c.Param(ParamAccount), Inactive) != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, &BoolResultResponse{Result: true})
}

func (a *Authentication) enable(c *gin.Context) {
	if Operator(a.mysql.Master, c.Param(ParamAccount), Active) != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, &BoolResultResponse{Result: true})
}
