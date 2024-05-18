package middleware

import (
	"go-server/internal/pkg/domains/models/dtos"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CheckRole() gin.HandlerFunc {
	return func(c *gin.Context) {
		isAdmin := c.MustGet("is_admin").(bool)
		if !isAdmin {
			c.JSON(http.StatusForbidden, dtos.BaseResponse{
				Code:    1,
				Message: "Forbidden",
				Error: &dtos.ErrorResponse{
					ErrorDetails: "Role invalid",
				},
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
