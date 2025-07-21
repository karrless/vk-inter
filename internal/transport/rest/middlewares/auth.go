package middlewares

import (
	"strings"
	"vk-inter/pkg/jwt"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware создает middleware для проверки JWT токена
func AuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// По умолчанию считаем неаутентифицированным
		c.Set("isAuthenticated", false)

		// Извлекаем токен из заголовка
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		tokenParts := strings.Split(authHeader, " ")

		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.Next()
			return
		}

		tokens, err := jwt.ValidateToken(tokenParts[1], secret)
		if err != nil {
			c.Next()
			return
		}

		// Токен валиден
		c.Set("isAuthenticated", true)
		c.Set("id", tokens.Subject)
		c.Next()
	}
}
