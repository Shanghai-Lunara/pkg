package qqbot

import (
	"fmt"
	"github.com/google/go-querystring/query"
	"strings"
)

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

// msg string

type Msg struct {
	str *strings.Builder
}

func NewMsg() *Msg {
	return &Msg{str: &strings.Builder{}}
}

func (m *Msg) ToString() string {
	return m.str.String()
}

// Text 文本
func (m *Msg) Text(text string) *Msg {
	// todo support color
	m.str.WriteString(text)
	return m
}

// Textln 文本
func (m *Msg) Textln(text string) *Msg {
	// todo support color
	m.str.WriteString(text)
	m.str.WriteString("\n")
	return m
}

// Text 文本
func (m *Msg) Textf(format string, a ...interface{}) *Msg {
	// todo support color
	m.str.WriteString(fmt.Sprintf(format, a...))
	return m
}

// At @某人
// qq @的 QQ 号, <=0 表示全体成员
func (m *Msg) At(qq int64) *Msg {
	// todo use cqcode struct
	str := ""
	if qq <= 0 {
		str = "[CQ:at,qq=all]"
	} else {
		str = fmt.Sprintf("[CQ:at,qq=%d]", qq)
	}
	m.str.WriteString(str)
	return m
}

// Link 分享链接
func (m *Msg) Link(url, title, content, image string) *Msg {
	// todo use cqcode struct
	m.str.WriteString(fmt.Sprintf("[CQ:share,url=%s,title=%s,content=%s,image=%s]", url, title, content, image))
	return m
}

// Face 表情
// id look https://github.com/kyubotics/coolq-http-api/wiki/%E8%A1%A8%E6%83%85-CQ-%E7%A0%81-ID-%E8%A1%A8
func (m *Msg) Face(id int) *Msg {
	m.str.WriteString(fmt.Sprintf("[CQ:face,id=%d]", id))
	return m
}
