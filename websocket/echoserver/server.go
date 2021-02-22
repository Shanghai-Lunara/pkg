package echoserver

import (
	"context"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/nevercase/lllidan/pkg/websocket"
	"k8s.io/klog/v2"
	"net/http"
)

type Server struct {
	members     *members
	connections *websocket.Connections
	server      *http.Server
}

func Init(ctx context.Context) *Server {
	s := &Server{
		members:     newMembers(ctx),
		connections: websocket.NewConnections(ctx),
	}
	router := gin.New()
	router.Use(cors.Default())
	router.GET("/", s.handler)
	server := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", 8081),
		Handler: router,
	}
	s.server = server
	go func() {
		if err := s.server.ListenAndServe(); err != nil {
			if err == http.ErrServerClosed {
				klog.Info("Server closed under request")
			} else {
				klog.V(2).Info("Server closed unexpected err:", err)
			}
		}
	}()
	return s
}

func (s *Server) handler(c *gin.Context) {
	s.connections.Handler(c.Writer, c.Request, s.members.newPlayer())
}

func (s *Server) Shutdown() {}
