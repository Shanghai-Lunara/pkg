package authentication

import (
	"github.com/Shanghai-Lunara/pkg/casbinrbac"
	"github.com/Shanghai-Lunara/pkg/jwttoken"
	"github.com/Shanghai-Lunara/pkg/zaplogger"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Authentication struct {
	mysql *casbinrbac.MysqlClusterPool
}

func New(router *gin.RouterGroup) *Authentication {
	authentication = &Authentication{
		mysql: casbinrbac.GetMysqlCluster(),
	}
	register(router)
	return authentication
}

var authentication *Authentication

const (
	AuthAccountLogin   = "/account/login"
	AuthAccountList    = "/account/list"
	AuthAccount        = "/account/:account/:pwd"
	AuthAccountDisable = "/account/:account"
)

const (
	ParamAccount  = "account"
	ParamPassword = "pwd"
)

func register(router *gin.RouterGroup) {
	router.Use(authentication.middleware())
	router.POST(AuthAccountLogin, authentication.login)
	router.GET(AuthAccountList, authentication.list)
	router.POST(AuthAccount, authentication.add)
	router.PUT(AuthAccount, authentication.reset)
	router.GET(AuthAccountDisable, authentication.disable)
}

func (a *Authentication) middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		switch c.FullPath() {
		case AuthAccountLogin:
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
			switch tokenClaims.Username {
			case "admin":
				c.Next()
			default:
				c.AbortWithStatus(http.StatusForbidden)
			}
		}
	}
}

func (a *Authentication) login(c *gin.Context) {
	req := &LoginRequest{}
	if c.ShouldBindJSON(req) != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	}
	acc, err := Query(a.mysql.Slave, req.Account)
	if err != nil {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	token, err := jwttoken.Generate(acc.Account)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, &LoginResponse{Token: token})
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
	req := AccountRequest{}
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
	req := AccountRequest{}
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
	req := &DisableRequest{}
	if c.ShouldBindJSON(req) != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if Disable(a.mysql.Master, req.Account) != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, &BoolResultResponse{Result: true})
}
