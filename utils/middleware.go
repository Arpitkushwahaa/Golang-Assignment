package utils

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// LoggingMiddleware logs HTTP requests
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(startTime)

		// Get status code
		statusCode := c.Writer.Status()

		// Log with structured fields
		logrus.WithFields(logrus.Fields{
			"status":     statusCode,
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"ip":         c.ClientIP(),
			"latency":    latency,
			"user_agent": c.Request.UserAgent(),
		}).Info("HTTP Request")

		// Log errors if any
		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				logrus.WithFields(logrus.Fields{
					"path":   c.Request.URL.Path,
					"method": c.Request.Method,
				}).Error(err.Error())
			}
		}
	}
}

// ErrorHandler handles panics and returns appropriate error responses
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logrus.WithFields(logrus.Fields{
					"error": err,
					"path":  c.Request.URL.Path,
				}).Error("Panic recovered")

				c.JSON(500, gin.H{
					"success": false,
					"error":   "Internal server error",
				})
			}
		}()
		c.Next()
	}
}
