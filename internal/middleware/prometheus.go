package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// HTTP Metrics
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	httpRequestErrors = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_request_errors_total",
			Help: "Total number of HTTP request errors",
		},
		[]string{"method", "path", "status"},
	)

	// Business Metrics
	MarksAddedTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "marks_added_total",
			Help: "Total number of marks added",
		},
	)

	AbsencesRecordedTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "absences_recorded_total",
			Help: "Total number of absences recorded",
		},
	)

	DocumentsGeneratedTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "documents_generated_total",
			Help: "Total number of documents generated",
		},
	)

	DocumentsGenerationErrors = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "documents_generation_errors_total",
			Help: "Total number of document generation errors",
		},
	)

	UsersLoginTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "users_login_total",
			Help: "Total number of user logins",
		},
	)

	UnauthorizedAccessAttempts = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "unauthorized_access_attempts_total",
			Help: "Total number of unauthorized access attempts",
		},
	)

	// WebSocket Metrics
	ActiveWebSocketConnections = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "active_websocket_connections",
			Help: "Number of active WebSocket connections",
		},
	)

	// Database Metrics
	DBQueryDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "db_query_duration_seconds",
			Help:    "Database query duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation"},
	)

	DBConnectionPoolSize = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "db_connection_pool_size",
			Help: "Number of connections in the database pool",
		},
	)

	SlowQueriesTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "db_slow_queries_total",
			Help: "Total number of slow database queries (>500ms)",
		},
	)
)

// PrometheusMiddleware records HTTP metrics
func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}

		c.Next()

		status := c.Writer.Status()
		duration := time.Since(start).Seconds()

		// Record metrics
		httpRequestsTotal.WithLabelValues(c.Request.Method, path, strconv.Itoa(status)).Inc()
		httpRequestDuration.WithLabelValues(c.Request.Method, path).Observe(duration)

		// Track errors (4xx and 5xx)
		if status >= 400 {
			httpRequestErrors.WithLabelValues(c.Request.Method, path, strconv.Itoa(status)).Inc()
		}
	}
}
