package middleware

import (
	"go-server/internal/pkg/domains/models/dtos"
	"go-server/internal/pkg/domains/models/entities"
	"go-server/pkg/shared/auth"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

const (
	CheckAuthenticationTokenNotSet = iota + 1
	CheckAuthenticationTokenInvalid
)

func CheckAuthentication(
	db *gorm.DB,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := c.Request.Header.Get("Authorization")
		if authorization == "" {
			c.JSON(http.StatusUnauthorized, dtos.BaseResponse{
				Code:    CheckAuthenticationTokenNotSet,
				Message: "Unauthorized",
				Error: &dtos.ErrorResponse{
					ErrorDetails: "Authorization token is not set",
				},
			})
			c.Abort()
			return
		}

		token := strings.Replace(authorization, "Bearer ", "", -1)
		fields := strings.Split(token, ".")

		decodedToken, err := auth.Decode(token)
		userID := decodedToken.Claims.(jwt.MapClaims)["user_id"]
		isAdmin := decodedToken.Claims.(jwt.MapClaims)["is_admin"]
		tokenID := decodedToken.Claims.(jwt.MapClaims)["token_id"]

		if len(fields) != 3 || !auth.VerifyJWT(token) || err != nil {
			c.JSON(http.StatusUnauthorized, dtos.BaseResponse{
				Code:    CheckAuthenticationTokenInvalid,
				Message: "Unauthorized",
				Error: &dtos.ErrorResponse{
					ErrorDetails: "Authorization token invalid",
				},
			})
			c.Abort()
			return
		}

		err = db.Where("token_id = ?", tokenID).First(&entities.UserToken{}).Error
		if err != nil {
			c.JSON(http.StatusUnauthorized, dtos.BaseResponse{
				Code:    CheckAuthenticationTokenInvalid,
				Message: "Unauthorized",
				Error: &dtos.ErrorResponse{
					ErrorDetails: "Authorization token invalid",
				},
			})
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		c.Set("is_admin", isAdmin)
		c.Set("token_id", tokenID)
		c.Next()
	}
}
