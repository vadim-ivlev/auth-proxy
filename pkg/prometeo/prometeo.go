package prometeo

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var httpRequestsTotal = promauto.NewCounter(prometheus.CounterOpts{
	Name: "http_requests_total",
	Help: "Тотальное количество запросов",
})

// CountersMiddleware считает запросы
func CountersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// считаем все кроме запросов на выдачу статистики
		if c.Request.URL.Path != "/metrics" {
			httpRequestsTotal.Inc()
		}
		c.Next()
	}
}
