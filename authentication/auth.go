package authentication

import (
	"github.com/Shanghai-Lunara/pkg/casbinrbac"
	"github.com/Shanghai-Lunara/pkg/jwttoken"
	"github.com/Shanghai-Lunara/pkg/zaplogger"
	"github.com/gin-gonic/gin"
)

type Authentication struct {
	mysql *casbinrbac.MysqlClusterPool
}

func New() *Authentication {
	authentication = &Authentication{
		mysql: casbinrbac.GetMysqlCluster(),
	}
	return authentication
}

var authentication *Authentication

const (
	AuthAccountLogin   = "/auth/account/login"
	AuthAccountList    = "/auth/account/list"
	AuthAccount        = "/auth/account/:account/:pwd"
	AuthAccountDisable = "/auth/account/:account"
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
		switch c.FullPath() {
		case AuthAccountLogin:
			c.Next()
		default:
			switch tokenClaims.Username {
			case "admin":
				c.Next()
			default:
				c.Abort()
			}
		}
	}
}

func (a *Authentication) login(c *gin.Context) {
	//a.mysql.Slave
}

func (a *Authentication) list(c *gin.Context) {
	//a.mysql.Master
}

func (a *Authentication) add(c *gin.Context) {

}

func (a *Authentication) reset(c *gin.Context) {

}

func (a *Authentication) disable(c *gin.Context) {

}
