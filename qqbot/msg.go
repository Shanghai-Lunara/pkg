package qqbot

import "github.com/google/go-querystring/query"

type Message interface {
	ToQueryString() string
}

type BotCallBack struct {
	Data    Data   `json:"data"`
	RetCode int32  `json:"retcode"`
	Status  string `json:"status"`
}

type Data struct {
	MessageId int32 `json:"message_id"`
}

type SendMsg struct {
	MessageType string `json:"message_type,omitempty" url:"message_type"`
	UserId      int64  `json:"user_id,omitempty" url:"user_id"`
	GroupId     int64  `json:"group_id,omitempty" url:"group_id"`
	Message     string `json:"message" url:"message"`
	AutoEscape  bool   `json:"auto_escape" url:"auto_escape"`
}

func (m SendMsg) ToQueryString() string {
	v, _ := query.Values(m)
	return v.Encode()
}
