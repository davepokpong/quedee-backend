package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthorizeAdminOrStaff() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _ := c.Get("role")
		if role != "admin" && role != "staff" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "only admin or staff can do this",
			})
			return
		}
		c.Next()
	}
}

func AuthorizeOnlyAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _ := c.Get("role")
		if role != "admin" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "only admin can do this",
			})
			return
		}
		c.Next()
	}
}
