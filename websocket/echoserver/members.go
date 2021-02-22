package echoserver

import (
	"context"
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

func (m *members) newPlayer() *handler {
	h := NewHandler(m.ctx, "", m.removedChan)
	return h
}
