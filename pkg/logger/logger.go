package logger

import (
	"context"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var globalLogger *logrus.Logger

// Init initializes the global logger
func Init() {
	globalLogger = logrus.New()
	
	// Set JSON formatter for structured logging
	globalLogger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "timestamp",
			logrus.FieldKeyLevel: "level",
			logrus.FieldKeyMsg:   "message",
		},
	})

	// Set log level based on environment
	level := os.Getenv("LOG_LEVEL")
	switch level {
	case "DEBUG":
		globalLogger.SetLevel(logrus.DebugLevel)
	case "INFO":
		globalLogger.SetLevel(logrus.InfoLevel)
	case "WARN":
		globalLogger.SetLevel(logrus.WarnLevel)
	case "ERROR":
		globalLogger.SetLevel(logrus.ErrorLevel)
	default:
		globalLogger.SetLevel(logrus.InfoLevel)
	}

	// Output to stdout
	globalLogger.SetOutput(os.Stdout)
}

// GetLogger returns the global logger instance
func GetLogger() *logrus.Logger {
	if globalLogger == nil {
		Init()
	}
	return globalLogger
}

// WithContext creates a logger with context information
func WithContext(ctx context.Context) *logrus.Entry {
	logger := GetLogger()
	entry := logger.WithContext(ctx)

	// Add correlation ID if available
	if correlationID := ctx.Value("correlation_id"); correlationID != nil {
		entry = entry.WithField("correlation_id", correlationID)
	}

	// Add trace ID if available
	if traceID := ctx.Value("trace_id"); traceID != nil {
		entry = entry.WithField("trace_id", traceID)
	}

	// Add tenant ID if available
	if tenantID := ctx.Value("tenant_id"); tenantID != nil {
		entry = entry.WithField("tenant_id", tenantID)
	}

	// Add user ID if available
	if userID := ctx.Value("user_id"); userID != nil {
		entry = entry.WithField("user_id", userID)
	}

	return entry
}

// WithField creates a logger with a single field
func WithField(key string, value interface{}) *logrus.Entry {
	return GetLogger().WithField(key, value)
}

// WithFields creates a logger with multiple fields
func WithFields(fields logrus.Fields) *logrus.Entry {
	return GetLogger().WithFields(fields)
}

// GinMiddleware returns a Gin middleware for structured logging
func GinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		duration := time.Since(start)
		status := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method
		userAgent := c.Request.UserAgent()

		fields := logrus.Fields{
			"status":      status,
			"duration_ms": duration.Milliseconds(),
			"client_ip":   clientIP,
			"method":      method,
			"path":        path,
			"user_agent":  userAgent,
		}

		if query != "" {
			fields["query"] = query
		}

		// Add tenant and user info if available
		if tenantID, exists := c.Get("tenant_id"); exists {
			fields["tenant_id"] = tenantID
		}

		if userID, exists := c.Get("user_id"); exists {
			fields["user_id"] = userID
		}

		// Add correlation ID if available
		if correlationID := c.GetHeader("X-Correlation-ID"); correlationID != "" {
			fields["correlation_id"] = correlationID
		}

		logger := GetLogger().WithFields(fields)

		switch {
		case status >= 500:
			logger.Error("HTTP request completed with server error")
		case status >= 400:
			logger.Warn("HTTP request completed with client error")
		case status >= 300:
			logger.Info("HTTP request completed with redirect")
		default:
			logger.Info("HTTP request completed successfully")
		}
	}
}

// Debug logs a debug message
func Debug(msg string, fields ...interface{}) {
	entry := GetLogger()
	for i := 0; i < len(fields); i += 2 {
		if i+1 < len(fields) {
			entry = entry.WithField(fields[i].(string), fields[i+1])
		}
	}
	entry.Debug(msg)
}

// Info logs an info message
func Info(msg string, fields ...interface{}) {
	entry := GetLogger()
	for i := 0; i < len(fields); i += 2 {
		if i+1 < len(fields) {
			entry = entry.WithField(fields[i].(string), fields[i+1])
		}
	}
	entry.Info(msg)
}

// Warn logs a warning message
func Warn(msg string, fields ...interface{}) {
	entry := GetLogger()
	for i := 0; i < len(fields); i += 2 {
		if i+1 < len(fields) {
			entry = entry.WithField(fields[i].(string), fields[i+1])
		}
	}
	entry.Warn(msg)
}

// Error logs an error message
func Error(msg string, fields ...interface{}) {
	entry := GetLogger()
	for i := 0; i < len(fields); i += 2 {
		if i+1 < len(fields) {
			entry = entry.WithField(fields[i].(string), fields[i+1])
		}
	}
	entry.Error(msg)
}

// Fatal logs a fatal message and exits
func Fatal(msg string, fields ...interface{}) {
	entry := GetLogger()
	for i := 0; i < len(fields); i += 2 {
		if i+1 < len(fields) {
			entry = entry.WithField(fields[i].(string), fields[i+1])
		}
	}
	entry.Fatal(msg)
}