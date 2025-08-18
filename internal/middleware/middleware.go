package middleware

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/time/rate"
)

// AuthMiddleware validates JWT tokens
func AuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Bearer token required"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Set("user_id", claims["user_id"])
			c.Set("tenant_id", claims["tenant_id"])
			c.Set("role", claims["role"])
		}

		c.Next()
	}
}

// TenantMiddleware extracts tenant information
func TenantMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID := c.GetHeader("X-Tenant-ID")
		if tenantID == "" {
			if claims, exists := c.Get("tenant_id"); exists {
				tenantID = claims.(string)
			}
		}

		if tenantID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Tenant ID required"})
			c.Abort()
			return
		}

		c.Set("tenant_id", tenantID)
		c.Next()
	}
}

// AdminMiddleware checks for admin role
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "Role not found"})
			c.Abort()
			return
		}

		if role != "admin" && role != "super_admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RateLimitConfig holds rate limiting configuration
type RateLimitConfig struct {
	Enabled bool `mapstructure:"enabled"`
	RPS     int  `mapstructure:"rps"`
	Burst   int  `mapstructure:"burst"`
}

// RateLimitMiddleware implements rate limiting per tenant
func RateLimitMiddleware(redisClient *redis.Client, config RateLimitConfig) gin.HandlerFunc {
	if !config.Enabled {
		return func(c *gin.Context) {
			c.Next()
		}
	}

	limiters := make(map[string]*rate.Limiter)

	return func(c *gin.Context) {
		tenantID, exists := c.Get("tenant_id")
		if !exists {
			tenantID = "default"
		}

		key := tenantID.(string)
		limiter, exists := limiters[key]
		if !exists {
			limiter = rate.NewLimiter(rate.Limit(config.RPS), config.Burst)
			limiters[key] = limiter
		}

		if !limiter.Allow() {
			c.Header("X-RateLimit-Limit", strconv.Itoa(config.RPS))
			c.Header("X-RateLimit-Remaining", "0")
			c.Header("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(time.Second).Unix(), 10))
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// CORSMiddleware handles CORS headers
func CORSMiddleware(allowedOrigins []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		
		// Check if origin is allowed
		allowed := false
		for _, allowedOrigin := range allowedOrigins {
			if allowedOrigin == "*" || allowedOrigin == origin {
				allowed = true
				break
			}
		}

		if allowed {
			c.Header("Access-Control-Allow-Origin", origin)
		}

		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-Tenant-ID")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// LoggingMiddleware logs HTTP requests
func LoggingMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return ""
	})
}

// SecurityMiddleware adds security headers
func SecurityMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		c.Header("Content-Security-Policy", "default-src 'self'")
		c.Next()
	}
}

// MetricsMiddleware for Prometheus metrics collection
func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path

		c.Next()

		status := c.Writer.Status()
		duration := time.Since(start)
		
		// Increment metrics (implementation depends on metrics package)
		// metrics.HTTPRequests.WithLabelValues(path, strconv.Itoa(status)).Inc()
		// metrics.HTTPDuration.WithLabelValues(path).Observe(duration.Seconds())
		
		_ = status
		_ = duration
		_ = path
	}
}