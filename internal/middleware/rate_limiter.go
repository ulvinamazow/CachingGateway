package middleware

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type Client struct {
	limiter *rate.Limiter
}

func RateLimiterMiddleware() gin.HandlerFunc {

	var (
		mu      sync.Mutex
		clients = make(map[string]*Client)
	)

	return func(c *gin.Context) {
		ip := c.ClientIP()

		mu.Lock()

		if _, exists := clients[ip]; !exists {
			clients[ip] = &Client{limiter: rate.NewLimiter(0.25, 15)}
		}

		cl := clients[ip]
		mu.Unlock()

		if !cl.limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "rate limit exceeded",
			})
			return
		}
		c.Next()
	}
}
