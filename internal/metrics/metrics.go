package metrics

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Metrics struct {
	httpRequestsTotal   *prometheus.CounterVec
	httpRequestDuration *prometheus.HistogramVec
	pastesCreatedTotal  prometheus.Counter
	dbQueriesTotal      *prometheus.CounterVec
}

func NewMetrics(reg prometheus.Registerer) *Metrics {
	m := &Metrics{
		httpRequestsTotal: promauto.With(reg).NewCounterVec(
			prometheus.CounterOpts{
				Name: "pastebin_http_requests_total",
				Help: "Total number of HTTP requests.",
			},
			[]string{"method", "path", "status_code"},
		),
		httpRequestDuration: promauto.With(reg).NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "pastebin_http_request_duration_seconds",
				Help:    "Duration of HTTP requests in seconds.",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "path"},
		),
		pastesCreatedTotal: promauto.With(reg).NewCounter(
			prometheus.CounterOpts{
				Name: "pastebin_pastes_created_total",
				Help: "Total number of created pastes.",
			},
		),
		dbQueriesTotal: promauto.With(reg).NewCounterVec(
			prometheus.CounterOpts{
				Name: "pastebin_db_queries_total",
				Help: "Total number of database queries executed.",
			},
			[]string{"operation"},
		),
	}
	return m
}

func (m *Metrics) IncPastesCreated() {
	m.pastesCreatedTotal.Inc()
}

func (m *Metrics) IncDBQuery(operation string) {
	m.dbQueriesTotal.WithLabelValues(operation).Inc()
}

func (m *Metrics) PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start).Seconds()
		statusCode := c.Writer.Status()
		path := c.FullPath()
		if path == "" {
			path = "not_found"
		}

		m.httpRequestDuration.WithLabelValues(c.Request.Method, path).Observe(duration)
		m.httpRequestsTotal.WithLabelValues(c.Request.Method, path, strconv.Itoa(statusCode)).Inc()
	}
}

func PrometheusHandler() gin.HandlerFunc {
	h := promhttp.Handler()
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
