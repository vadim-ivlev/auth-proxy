package prometeo

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var HttpRequestsTotal = promauto.NewCounter(prometheus.CounterOpts{
	Name: "http_requests_total",
	Help: "Тотальное количество запросов",
})

var GraphQLErrorsTotal = promauto.NewCounter(prometheus.CounterOpts{
	Name: "graphql_errors_total",
	Help: "Тотальное количество ошибок GraphQL",
})

// CountersMiddleware считает запросы
func CountersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// считаем все кроме запросов на выдачу статистики
		if c.Request.URL.Path != "/metrics" {
			HttpRequestsTotal.Inc()
		}
		c.Next()
	}
}
