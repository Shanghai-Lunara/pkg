package echoserver

import (
	"context"
	"sync"
)

type handler struct {
	sync.Once
	id         int32
	router     string
	removeChan chan<- int32
	outputChan chan<- []byte
	pingFunc   func()
	closeFunc  func()
	ctx        context.Context
	cancel     context.CancelFunc
}

func NewHandler(ctx context.Context, router string, removeChan chan<- int32) *handler {
	sub, cancel := context.WithCancel(ctx)
	c := &handler{
		id:         0,
		router:     router,
		removeChan: removeChan,
		ctx:        sub,
		cancel:     cancel,
	}
	return c
}

func (h *handler) RegisterId(id int32) {
	h.id = id
}

func (h *handler) RegisterRemoveChan(ch chan<- int32) {
	h.removeChan = ch
}

func (h *handler) RegisterConnWriteChan(ch chan<- []byte) {
	h.outputChan = ch
}

func (h *handler) RegisterConnClose(do func()) {
	h.closeFunc = do
}

func (h *handler) RegisterConnPing(do func()) {
	h.pingFunc = do
}

func (h *handler) Handler(in []byte) (res []byte, err error) {
	h.pingFunc()
	return in, nil
}

func (h *handler) Run() {}

func (h *handler) Close() {
	h.Once.Do(func() {
		h.removeChan <- h.id
		h.cancel()
	})
}
