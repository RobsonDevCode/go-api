package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/RobsonDevCode/go-profile-service/src/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

func JWTAuthMiddleWare(config *config.Config, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User Authorized"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User Authorized"})
			return
		}

		jwtString := parts[1]

		if config.JWTSettings.Key == "" {
			logger.Sugar().Error("config error jwt key not found")
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong please try again later!"})
			return
		}

		token, err := jwt.Parse(jwtString, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(config.JWTSettings.Key), nil
		})

		if err != nil {
			if err == jwt.ErrTokenSignatureInvalid {
				logger.Sugar().Errorf("Invalid token signature: %w", err)
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User does not have access "})
				return
			}

			logger.Sugar().Errorf("error, %w", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User unauthorized"})
			return
		}

		if !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User unauthorized"})
			return
		}
	}
}
