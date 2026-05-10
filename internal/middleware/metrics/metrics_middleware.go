package metrics

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	httpRequestTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of http request.",
		},

		[]string{"method", "endpoint", "status"},
	)

	httpRequestsInFlight = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_requests_in_flight",
			Help: "Curre number of in-flight requests.",
		},
	)

	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests",
			Buckets: []float64{0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1},
		},

		[]string{"method", "path", "status"},
	)
)

func RegisterCollectors() {
	prometheus.MustRegister(httpRequestTotal)
	prometheus.MustRegister(httpRequestsInFlight)
	prometheus.MustRegister(httpRequestDuration)

}

func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		start := time.Now()

		path := c.FullPath()

		if path == "" {
			path = c.Request.URL.Path
		}

		httpRequestsInFlight.Inc()
		defer httpRequestsInFlight.Dec()

		c.Next()

		status := fmt.Sprintf("%d", c.Writer.Status())
		method := c.Request.Method
		duration := time.Since(start).Seconds()

		httpRequestTotal.WithLabelValues(method, path, status).Inc()

		httpRequestDuration.WithLabelValues(method, path, status).Observe(duration)
	}
}
