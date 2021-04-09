package casbinrbac

import (
	"fmt"
	"github.com/Shanghai-Lunara/pkg/zaplogger"
)

const (
	CodeSuccess = iota
	CodeError
)

type BoolResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Result  bool   `json:"result"`
}

//GetPolicy	{"policy": [["p2","namespace2","data2","read"],["p1","namespace1","data1","read"],["p1","namespace3","data3","write"]]}
//GetGroupingPolicy	{"users": [["alice","p2","namespace2"],["alice","p1","namespace1"]]}

type Policy struct {
	Role      string `json:"role"`
	Namespace string `json:"namespace"`
	Object    string `json:"object"`
	Action    string `json:"action"`
}

func ConvertPolicy(in []string) Policy {
	if len(in) != 4 {
		zaplogger.Sugar().Fatal(fmt.Errorf("invalid policy size: expected 4, got %d, pvals: %#v", len(in), in))
	}
	return Policy{
		Role:      in[0],
		Namespace: in[1],
		Object:    in[2],
		Action:    in[3],
	}
}

type ListPolicyResponse struct {
	Code     int      `json:"code"`
	Message  string   `json:"message"`
	Policies []Policy `json:"policies"`
}

type GroupingPolicy struct {
	User      string `json:"user"`
	Role      string `json:"role"`
	Namespace string `json:"namespace"`
}

func ConvertGroupingPolicy(in []string) GroupingPolicy {
	if len(in) != 3 {
		zaplogger.Sugar().Fatal(fmt.Errorf("invalid GroupingPolicy size: expected 3, got %d, pvals: %#v", len(in), in))
	}
	return GroupingPolicy{
		User:      in[0],
		Role:      in[1],
		Namespace: in[2],
	}
}

type ListGroupingPolicyResponse struct {
	Code             int              `json:"code"`
	Message          string           `json:"message"`
	GroupingPolicies []GroupingPolicy `json:"grouping_policies"`
}
