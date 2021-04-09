package casbinrbac

import (
	"github.com/Shanghai-Lunara/pkg/zaplogger"
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/gin-gonic/gin"
	"net/http"
)

type RBAC struct {
	e *casbin.Enforcer
}

var rbac *RBAC

func New(rulePath string, dsn string, router *gin.RouterGroup) *RBAC {
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
	register(router)
	return rbac
}

const (
	AddPermissionForRole    = "/casbin/permission/add/:role/:namespace/:permission/:action"
	DeletePermissionForRole = "/casbin/permission/delete/:role/:namespace/:permission/:action"
	AddRoleForUser          = "/casbin/role/add/:user/:namespace/:role"
	DeleteRoleForUser       = "/casbin/role/delete/:user/:namespace/:role"
	ListPolicy              = "/casbin/policy/list"
	ListGroupingPolicy      = "/casbin/groupingpolicy/list"
	FilteredGroupingPolicy  = "/casbin/groupingpolicy/filter"
)

func register(router *gin.RouterGroup) {
	router.Use(rbac.auth())
	router.GET(AddPermissionForRole, rbac.AddPermissionForRole)
	router.GET(DeletePermissionForRole, rbac.DeletePermissionForRole)
	router.GET(AddRoleForUser, rbac.AddRoleForUser)
	router.GET(DeleteRoleForUser, rbac.DeleteRoleForUser)
	router.GET(ListPolicy, rbac.ListPolicy)
	router.GET(ListGroupingPolicy, rbac.ListGroupingPolicy)
	router.GET(FilteredGroupingPolicy, rbac.ListGroupingPolicy)
}

func (r *RBAC) auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := "admin"
		switch c.FullPath() {
		case FilteredGroupingPolicy:
			c.Next()
		default:
			switch user {
			case "admin":
				c.Next()
			default:
				c.Abort()
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

func (r *RBAC) AddPermissionForRole(c *gin.Context) {
	ok, err := r.e.AddPermissionForUser(c.Param("role"), c.Param("namespace"), c.Param("permission"), c.Param("action"))
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

func (r *RBAC) DeletePermissionForRole(c *gin.Context) {
	ok, err := r.e.DeletePermissionForUser(c.Param("role"), c.Param("namespace"), c.Param("permission"), c.Param("action"))
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

func (r *RBAC) AddRoleForUser(c *gin.Context) {
	ok, err := r.e.AddRoleForUser(c.Param("user"), c.Param("role"), c.Param("namespace"))
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

func (r *RBAC) DeleteRoleForUser(c *gin.Context) {
	ok, err := r.e.DeleteRoleForUser(c.Param("user"), c.Param("role"), c.Param("namespace"))
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

func (r *RBAC) ListPolicy(c *gin.Context) {
	t := r.e.GetPolicy()
	p := make([]Policy, 0)
	for _, v := range t {
		p = append(p, ConvertPolicy(v))
	}
	c.JSON(http.StatusOK, ListPolicyResponse{
		Code:     0,
		Message:  "",
		Policies: p,
	})
}

func (r *RBAC) ListGroupingPolicy(c *gin.Context) {
	t := r.e.GetGroupingPolicy()
	p := make([]GroupingPolicy, 0)
	for _, v := range t {
		p = append(p, ConvertGroupingPolicy(v))
	}
	c.JSON(http.StatusOK, ListGroupingPolicyResponse{
		Code:             0,
		Message:          "",
		GroupingPolicies: p,
	})
}

func (r *RBAC) FilteredGroupingPolicy(c *gin.Context) {
	t := r.e.GetFilteredGroupingPolicy(0, "user")
	p := make([]GroupingPolicy, 0)
	for _, v := range t {
		p = append(p, ConvertGroupingPolicy(v))
	}
	c.JSON(http.StatusOK, ListGroupingPolicyResponse{
		Code:             0,
		Message:          "",
		GroupingPolicies: p,
	})
}

func (r *RBAC) save() {
	if err := r.e.SavePolicy(); err != nil {
		zaplogger.Sugar().Error(err)
	}
}

// "alice", "namespace1",  "data1", "read"
func Enforce(userOrRole string, namespace string, object string, action string) (bool, error) {
	if rbac == nil {
		zaplogger.Sugar().Fatal("error: nil RBAC, please call New() before Enforce()")
	}
	return rbac.e.Enforce(userOrRole, namespace, object, action)
}
