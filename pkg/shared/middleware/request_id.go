package middleware

import (
	"go-server/pkg/shared/logging/hooks"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RequestID adds a unique request ID to the context and response headers.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := uuid.New().String()

		// Add the request ID to the context.
		c.Set(hooks.RequestIDField, requestID)

		// Add the request ID to the response headers.
		c.Writer.Header().Set("X-Request-ID", requestID)

		c.Next()
	}
}
