package casbinrbac

import (
	"fmt"
	"github.com/Shanghai-Lunara/pkg/jwttoken"
	"github.com/Shanghai-Lunara/pkg/zaplogger"
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
)

type RBAC struct {
	mu           sync.RWMutex
	relativePath string
	e            *casbin.Enforcer
}

var rbac *RBAC

func NewWithMysqlConf(rulePath string, mysqlConfPath string, relativePath string, router *gin.Engine) *RBAC {
	LoadMysqlConf(mysqlConfPath)
	a, err := gormadapter.NewAdapter("mysql", MasterDsn())
	if err != nil {
		zaplogger.Sugar().Fatal(err)
	}
	e, err := casbin.NewEnforcer(rulePath, a)
	if err != nil {
		zaplogger.Sugar().Fatal(err)
	}
	if err = e.LoadPolicy(); err != nil {
		zaplogger.Sugar().Fatal(err)
	}
	rbac = &RBAC{
		relativePath: relativePath,
		e:            e,
	}
	register(router.Group(relativePath))
	return rbac
}

func NewWithDsnString(rulePath string, dsn string, relativePath string, router *gin.Engine) *RBAC {
	a, err := gormadapter.NewAdapter("mysql", dsn)
	if err != nil {
		zaplogger.Sugar().Fatal(err)
	}
	e, err := casbin.NewEnforcer(rulePath, a)
	if err != nil {
		zaplogger.Sugar().Fatal(err)
	}
	if err = e.LoadPolicy(); err != nil {
		zaplogger.Sugar().Fatal(err)
	}
	rbac = &RBAC{
		e: e,
	}
	register(router.Group(relativePath))
	return rbac
}

const (
	RouterAddPermissionForRole    = "/casbin/permission/add/:role/:namespace/:permission/:action"
	RouterDeletePermissionForRole = "/casbin/permission/delete/:role/:namespace/:permission/:action"
	RouterAddRoleForUser          = "/casbin/role/add/:user/:namespace/:role"
	RouterDeleteRoleForUser       = "/casbin/role/delete/:user/:namespace/:role"
	RouterListPolicy              = "/casbin/policy/list"
	RouterListGroupingPolicy      = "/casbin/groupingpolicy/list"
	RouterFilterGroupingPolicy    = "/casbin/groupingpolicy/filter"
)

func register(router *gin.RouterGroup) {
	router.Use(rbac.auth())
	router.GET(RouterAddPermissionForRole, rbac.AddPermissionForRoleHandler)
	router.GET(RouterDeletePermissionForRole, rbac.DeletePermissionForRoleHandler)
	router.GET(RouterAddRoleForUser, rbac.AddRoleForUserHandler)
	router.GET(RouterDeleteRoleForUser, rbac.DeleteRoleForUserHandler)
	router.GET(RouterListPolicy, rbac.ListPolicyHandler)
	router.GET(RouterListGroupingPolicy, rbac.ListGroupingPolicyHandler)
	router.GET(RouterFilterGroupingPolicy, rbac.FilterGroupingPolicyHandler)
}

func (r *RBAC) auth() gin.HandlerFunc {
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
		switch c.FullPath() {
		case FullPath(r.relativePath, RouterFilterGroupingPolicy):
			c.Next()
		default:
			switch tokenClaims.IsAdmin {
			case true:
				c.Next()
			default:
				c.AbortWithStatus(http.StatusForbidden)
			}
		}
	}
}

func (r *RBAC) boolResponse(c *gin.Context, code int, ok bool, msg string) {
	c.JSON(http.StatusOK, BoolResponse{
		Code:    code,
		Message: msg,
		Result:  ok,
	})
}

func (r *RBAC) AddPermissionForRoleHandler(c *gin.Context) {
	r.mu.Lock()
	defer r.mu.Unlock()
	ok, err := r.e.AddPermissionForUser(c.Param(Role), c.Param(Namespace), c.Param(Permission), c.Param(Action))
	if err != nil {
		zaplogger.Sugar().Error(err)
		r.boolResponse(c, CodeError, false, err.Error())
		return
	}
	if ok {
		r.save()
	}
	r.boolResponse(c, CodeSuccess, ok, "")
}

func (r *RBAC) DeletePermissionForRoleHandler(c *gin.Context) {
	r.mu.Lock()
	defer r.mu.Unlock()
	ok, err := r.e.DeletePermissionForUser(c.Param(Role), c.Param(Namespace), c.Param(Permission), c.Param(Action))
	if err != nil {
		zaplogger.Sugar().Error(err)
		r.boolResponse(c, CodeError, false, err.Error())
		return
	}
	if ok {
		r.save()
	}
	r.boolResponse(c, CodeSuccess, ok, "")
}

func (r *RBAC) AddRoleForUserHandler(c *gin.Context) {
	r.mu.Lock()
	defer r.mu.Unlock()
	ok, err := r.e.AddRoleForUser(c.Param(User), c.Param(Role), c.Param(Namespace))
	if err != nil {
		zaplogger.Sugar().Error(err)
		r.boolResponse(c, CodeError, false, err.Error())
		return
	}
	if ok {
		r.save()
	}
	r.boolResponse(c, CodeSuccess, ok, "")
}

func (r *RBAC) DeleteRoleForUserHandler(c *gin.Context) {
	r.mu.Lock()
	defer r.mu.Unlock()
	ok, err := r.e.DeleteRoleForUser(c.Param(User), c.Param(Role), c.Param(Namespace))
	if err != nil {
		zaplogger.Sugar().Error(err)
		r.boolResponse(c, CodeError, false, err.Error())
		return
	}
	if ok {
		r.save()
	}
	r.boolResponse(c, CodeSuccess, ok, "")
}

func ListPolicy() []Policy {
	rbac.mu.RLock()
	defer rbac.mu.RUnlock()
	t := rbac.e.GetPolicy()
	p := make([]Policy, 0)
	for _, v := range t {
		p = append(p, ConvertPolicy(v))
	}
	return p
}

func (r *RBAC) ListPolicyHandler(c *gin.Context) {
	rbac.mu.RLock()
	defer rbac.mu.RUnlock()
	t := r.e.GetPolicy()
	p := make([]Policy, 0)
	for _, v := range t {
		p = append(p, ConvertPolicy(v))
	}
	c.JSON(http.StatusOK, ListPolicyResponse{
		Code:     0,
		Message:  "",
		Policies: ListPolicy(),
	})
}

func ListGroupingPolicy() []GroupingPolicy {
	rbac.mu.RLock()
	defer rbac.mu.RUnlock()
	t := rbac.e.GetGroupingPolicy()
	p := make([]GroupingPolicy, 0)
	for _, v := range t {
		p = append(p, ConvertGroupingPolicy(v))
	}
	return p
}

func (r *RBAC) ListGroupingPolicyHandler(c *gin.Context) {
	c.JSON(http.StatusOK, ListGroupingPolicyResponse{
		Code:             0,
		Message:          "",
		GroupingPolicies: ListGroupingPolicy(),
	})
}

func ListFilterGroupingPolicy(username string) []GroupingPolicy {
	rbac.mu.RLock()
	defer rbac.mu.RUnlock()
	t := rbac.e.GetFilteredGroupingPolicy(0, username)
	p := make([]GroupingPolicy, 0)
	for _, v := range t {
		p = append(p, ConvertGroupingPolicy(v))
	}
	return p
}

func (r *RBAC) FilterGroupingPolicyHandler(c *gin.Context) {
	token := c.Request.Header.Get(jwttoken.TokenKey)
	if token == "" {
		return
	}
	tokenClaims, err := jwttoken.Parse(token)
	if err != nil {
		zaplogger.Sugar().Error(err)
		return
	}
	c.JSON(http.StatusOK, ListGroupingPolicyResponse{
		Code:             0,
		Message:          "",
		GroupingPolicies: ListFilterGroupingPolicy(tokenClaims.Username),
	})
}

func (r *RBAC) save() {
	if err := r.e.SavePolicy(); err != nil {
		zaplogger.Sugar().Error(err)
	}
}

// "alice", "namespace1",  "data1", "read"
func Enforce(userOrRole string, namespace string, object string, action string) (bool, error) {
	rbac.mu.RLock()
	defer rbac.mu.RUnlock()
	if rbac == nil {
		zaplogger.Sugar().Fatal("error: nil RBAC, please call New() before Enforce()")
	}
	return rbac.e.Enforce(userOrRole, namespace, object, action)
}

func ListPoliciesByUsername(username string) []Policy {
	//rbac.mu.RLock()
	//defer rbac.mu.RUnlock()
	p := make([]Policy, 0)
	gp := ListFilterGroupingPolicy(username)
	if len(gp) == 0 {
		return p
	}
	tmp := make(map[string]bool, 0)
	for _, v := range gp {
		tmp[fmt.Sprintf("%s:%s", v.Role, v.Namespace)] = true
	}
	policies := ListPolicy()
	for _, v := range policies {
		if _, ok := tmp[fmt.Sprintf("%s:%s", v.Role, v.Namespace)]; ok {
			p = append(p, v)
		}
	}
	return p
}
