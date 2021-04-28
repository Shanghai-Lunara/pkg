package authentication

type Status int

const (
	Active   Status = 0
	Inactive Status = 1
)

type Account struct {
	Account    string `json:"account"`
	Password   string `json:"password"`
	CreateTime int    `json:"create_time"`
	Status     Status `json:"status"`
}

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ListResponse struct {
	Response
	Items []Account `json:"items"`
}

type BoolResultResponse struct {
	Response
	Result bool `json:"result"`
}

type LoginRequest struct {
	Account  string `json:"account"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Response
	Token string `json:"token"`
}

type ResetRequest struct {
	Account  string `json:"account"`
	Password string `json:"password"`
}

type ResetResponse struct {
	BoolResultResponse
}
