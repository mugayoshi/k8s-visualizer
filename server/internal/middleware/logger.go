// internal/middleware/logger.go
package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger is a custom logging middleware
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Get status code
		statusCode := c.Writer.Status()

		// Get client IP
		clientIP := c.ClientIP()

		// Get request method
		method := c.Request.Method

		// Build query string
		if raw != "" {
			path = path + "?" + raw
		}

		// Log format: [timestamp] status_code | latency | client_ip | method | path | error
		errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()

		log.Printf("[%s] %d | %13v | %15s | %-7s %s | %s",
			start.Format("2006/01/02 - 15:04:05"),
			statusCode,
			latency,
			clientIP,
			method,
			path,
			errorMessage,
		)
	}
}

// LoggerWithConfig allows customization
type LoggerConfig struct {
	// SkipPaths is a list of paths to skip logging
	SkipPaths []string
	// SkipPathRegexps is a list of regex patterns to skip logging
	SkipPathRegexps []string
}

func LoggerWithConfig(config LoggerConfig) gin.HandlerFunc {
	skipPaths := make(map[string]bool, len(config.SkipPaths))
	for _, path := range config.SkipPaths {
		skipPaths[path] = true
	}

	return func(c *gin.Context) {
		path := c.Request.URL.Path

		// Skip logging for certain paths (like /health)
		if skipPaths[path] {
			c.Next()
			return
		}

		start := time.Now()
		raw := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method

		if raw != "" {
			path = path + "?" + raw
		}

		errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()

		log.Printf("[%s] %d | %13v | %15s | %-7s %s | %s",
			start.Format("2006/01/02 - 15:04:05"),
			statusCode,
			latency,
			clientIP,
			method,
			path,
			errorMessage,
		)
	}
}
