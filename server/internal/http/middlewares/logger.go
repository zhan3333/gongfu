package middlewares

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
)

func JSONLogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()

		// Process Request
		c.Next()

		if strings.HasPrefix(c.Request.URL.Path, "/web") {
			return
		}
		if c.Request.URL.Path == "/favicon.ico" {
			return
		}

		// Stop timer
		duration := time.Since(start)

		entry := log.WithFields(log.Fields{
			"client_ip":  c.RemoteIP(),
			"duration":   duration.String(),
			"method":     c.Request.Method,
			"path":       c.Request.RequestURI,
			"status":     c.Writer.Status(),
			"referrer":   c.Request.Referer(),
			"request_id": c.Writer.Header().Get("Request-Id"),
			// "api_version": util.ApiVersion,
		})

		if c.Writer.Status() >= 500 {
			entry.Error(c.Errors.String())
		} else {
			entry.Info("access_log")
		}
	}
}
