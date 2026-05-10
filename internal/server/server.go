package server

import (
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/ulvinamazow/caching-gateway/internal/middleware"
	"github.com/ulvinamazow/caching-gateway/internal/middleware/metrics"
	"github.com/ulvinamazow/caching-gateway/internal/proxy"
)

type Server struct {
	port    string
	handler *proxy.Handler
}

func NewServer(port string, handler *proxy.Handler) *Server {
	return &Server{
		port:    port,
		handler: handler,
	}
}

func (server *Server) Start() error {
	router := gin.Default()

	router.Use(middleware.RateLimiterMiddleware())
	router.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		AllowMethods:     []string{"GET", "POST", "DELETE"},
		AllowHeaders:     []string{"Content-Type", "Authorization", "Origin"},
		ExposeHeaders:    []string{"X-Cache", "Content-Length"},
		AllowCredentials: true,
	}))
	router.Use(metrics.MetricsMiddleware())

	router.Any("/*path", func(c *gin.Context) {
		switch c.Param("path") {
		case "/health":
			server.handler.HealthCheck(c)
		case "/metrics":
			gin.WrapH(promhttp.Handler())(c)
		default:
			server.handler.Handle(c)
		}
	})

	addr := fmt.Sprintf(":%s", server.port)

	return router.Run(addr)
}
