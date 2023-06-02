package middlewares

import (
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func ValidateToken(accessToken string) (*jwt.Token, error) {
	return jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid token")
		}
		return []byte("secret"), nil
	})
}

func AuthorizeJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		const BEARER_SCHEMA = "Bearer "
		authHeader := c.GetHeader("Authorization")
		if len(authHeader) == 0 {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		tokenString := authHeader[len(BEARER_SCHEMA):]
		token, err := ValidateToken(tokenString)
		if token.Valid {
			claims := token.Claims.(jwt.MapClaims)
			c.Set("username", claims["username"])
			c.Set("role", claims["role"])
			// c.Set("user", claims["user"])
		} else {
			fmt.Println(err)
			c.AbortWithStatus(http.StatusUnauthorized)
		}
	}
}
