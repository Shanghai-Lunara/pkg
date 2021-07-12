package qqbot

import (
	"encoding/json"
	"fmt"
	"github.com/Shanghai-Lunara/pkg/zaplogger"
	"io/ioutil"
	"net/http"
)

type Client struct {
	addr string
	port int
}

func NewClient(addr string, port int) Client {
	return Client{
		addr: addr,
		port: port,
	}
}

func (c Client) SendGroupMsg(group int64, content string) (*BotCallBack, error) {
	msg := SendMsg{
		MessageType: "group",
		GroupId:     group,
		Message:     content,
		AutoEscape:  false,
	}
	return c.send(ApiSendMsg, msg)
}

func (c Client) SendPanicMsg(group int64, project, server, errMsg, stackTrace string) {
	c.SendGroupMsg(group, fmt.Sprintf(TmpErrMsg, project, server, errMsg, stackTrace))
}

func (c *Client) send(api string, msg Message) (*BotCallBack, error) {

	url := fmt.Sprintf("%s:%d/%s?%s", c.addr, c.port, api, msg.ToQueryString())
	method := "GET"
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		zaplogger.Sugar().Error(err)
		return nil, err
	}
	res, err := client.Do(req)
	if err != nil {
		zaplogger.Sugar().Error(err)
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		zaplogger.Sugar().Error(err)
		return nil, err
	}
	cb := &BotCallBack{}
	err1 := json.Unmarshal(body, cb)
	if err1 != nil {
		zaplogger.Sugar().Error("unmarshal fail", "err", err1.Error(), "data", string(body))
		return nil, err1
	}
	return cb, nil
}
