package metrics

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// HTTP Metrics
	HTTPRequests = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "{{MCP_NAME}}_http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status", "tenant_id"},
	)

	HTTPDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "{{MCP_NAME}}_http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint", "tenant_id"},
	)

	HTTPRequestsInProgress = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "{{MCP_NAME}}_http_requests_in_progress",
			Help: "Number of HTTP requests currently in progress",
		},
		[]string{"method", "endpoint"},
	)

	// Database Metrics
	DatabaseConnections = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "{{MCP_NAME}}_database_connections",
			Help: "Number of database connections",
		},
		[]string{"database", "state"},
	)

	DatabaseOperations = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "{{MCP_NAME}}_database_operations_total",
			Help: "Total number of database operations",
		},
		[]string{"database", "operation", "status"},
	)

	DatabaseQueryDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "{{MCP_NAME}}_database_query_duration_seconds",
			Help:    "Database query duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"database", "operation"},
	)

	// Redis Metrics
	RedisOperations = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "{{MCP_NAME}}_redis_operations_total",
			Help: "Total number of Redis operations",
		},
		[]string{"operation", "status"},
	)

	RedisConnectionsActive = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "{{MCP_NAME}}_redis_connections_active",
			Help: "Number of active Redis connections",
		},
	)

	// NATS Metrics
	NATSMessages = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "{{MCP_NAME}}_nats_messages_total",
			Help: "Total number of NATS messages",
		},
		[]string{"subject", "type", "status"},
	)

	NATSMessageDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "{{MCP_NAME}}_nats_message_duration_seconds",
			Help:    "NATS message processing duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"subject", "type"},
	)

	// AI Service Metrics
	AIOperations = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "{{MCP_NAME}}_ai_operations_total",
			Help: "Total number of AI operations",
		},
		[]string{"service", "operation", "status"},
	)

	AIOperationDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "{{MCP_NAME}}_ai_operation_duration_seconds",
			Help:    "AI operation duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"service", "operation"},
	)

	AITokensUsed = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "{{MCP_NAME}}_ai_tokens_used_total",
			Help: "Total number of AI tokens used",
		},
		[]string{"service", "type"},
	)

	// Business Logic Metrics (to be customized per MCP)
	{{BUSINESS_METRIC_1}} = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "{{MCP_NAME}}_{{METRIC_1_NAME}}_total",
			Help: "{{METRIC_1_DESCRIPTION}}",
		},
		[]string{"{{METRIC_1_LABELS}}"},
	)

	{{BUSINESS_METRIC_2}} = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "{{MCP_NAME}}_{{METRIC_2_NAME}}",
			Help: "{{METRIC_2_DESCRIPTION}}",
		},
		[]string{"{{METRIC_2_LABELS}}"},
	)

	{{BUSINESS_METRIC_3}} = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "{{MCP_NAME}}_{{METRIC_3_NAME}}_duration_seconds",
			Help:    "{{METRIC_3_DESCRIPTION}}",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"{{METRIC_3_LABELS}}"},
	)
)

// Init initializes the metrics
func Init() {
	// Register custom metrics if needed
	// This function can be extended to set up custom metrics collectors
}

// GinMiddleware returns a Gin middleware for metrics collection
func GinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		method := c.Request.Method
		endpoint := c.FullPath()
		
		// Handle cases where FullPath() returns empty string
		if endpoint == "" {
			endpoint = c.Request.URL.Path
		}

		// Get tenant ID for labeling
		tenantID := "unknown"
		if tid, exists := c.Get("tenant_id"); exists {
			tenantID = tid.(string)
		}

		// Increment in-progress requests
		HTTPRequestsInProgress.WithLabelValues(method, endpoint).Inc()

		// Process request
		c.Next()

		// Calculate duration
		duration := time.Since(start)
		status := strconv.Itoa(c.Writer.Status())

		// Record metrics
		HTTPRequests.WithLabelValues(method, endpoint, status, tenantID).Inc()
		HTTPDuration.WithLabelValues(method, endpoint, tenantID).Observe(duration.Seconds())
		HTTPRequestsInProgress.WithLabelValues(method, endpoint).Dec()
	}
}

// RecordDatabaseOperation records a database operation metric
func RecordDatabaseOperation(database, operation, status string, duration time.Duration) {
	DatabaseOperations.WithLabelValues(database, operation, status).Inc()
	DatabaseQueryDuration.WithLabelValues(database, operation).Observe(duration.Seconds())
}

// RecordRedisOperation records a Redis operation metric
func RecordRedisOperation(operation, status string) {
	RedisOperations.WithLabelValues(operation, status).Inc()
}

// RecordNATSMessage records a NATS message metric
func RecordNATSMessage(subject, messageType, status string, duration time.Duration) {
	NATSMessages.WithLabelValues(subject, messageType, status).Inc()
	NATSMessageDuration.WithLabelValues(subject, messageType).Observe(duration.Seconds())
}

// RecordAIOperation records an AI operation metric
func RecordAIOperation(service, operation, status string, duration time.Duration, tokensUsed int) {
	AIOperations.WithLabelValues(service, operation, status).Inc()
	AIOperationDuration.WithLabelValues(service, operation).Observe(duration.Seconds())
	if tokensUsed > 0 {
		AITokensUsed.WithLabelValues(service, "total").Add(float64(tokensUsed))
	}
}

// SetDatabaseConnections sets the current number of database connections
func SetDatabaseConnections(database, state string, count int) {
	DatabaseConnections.WithLabelValues(database, state).Set(float64(count))
}

// SetRedisActiveConnections sets the current number of active Redis connections
func SetRedisActiveConnections(count int) {
	RedisConnectionsActive.Set(float64(count))
}

// Custom business metric helpers (to be customized per MCP)

// Record{{BUSINESS_OPERATION_1}} records {{BUSINESS_OPERATION_1_DESCRIPTION}}
func Record{{BUSINESS_OPERATION_1}}(labels ...string) {
	{{BUSINESS_METRIC_1}}.WithLabelValues(labels...).Inc()
}

// Set{{BUSINESS_OPERATION_2}} sets {{BUSINESS_OPERATION_2_DESCRIPTION}}
func Set{{BUSINESS_OPERATION_2}}(value float64, labels ...string) {
	{{BUSINESS_METRIC_2}}.WithLabelValues(labels...).Set(value)
}

// Record{{BUSINESS_OPERATION_3}} records {{BUSINESS_OPERATION_3_DESCRIPTION}}
func Record{{BUSINESS_OPERATION_3}}(duration time.Duration, labels ...string) {
	{{BUSINESS_METRIC_3}}.WithLabelValues(labels...).Observe(duration.Seconds())
}