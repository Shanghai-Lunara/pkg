package echoserver

import (
	"context"
	"github.com/nevercase/lllidan/pkg/websocket/handler"
)

type members struct {
	removedChan chan int32
	ctx         context.Context
}

func newMembers(ctx context.Context) *members {
	m := &members{
		removedChan: make(chan int32, 1024),
		ctx:         ctx,
	}
	go m.loop()
	return m
}

func (m *members) loop() {
	for {
		select {
		case _, isClose := <-m.removedChan:
			if !isClose {
				return
			}
		}
	}
}

func (m *members) newPlayer() handler.Interface {
	h := handler.NewHandler(m.ctx, "", m.handler, m.removedChan)
	return h
}

func (m *members) handler(in []byte, handler handler.Interface) (res []byte, err error) {
	handler.Ping()
	return in, nil
}
