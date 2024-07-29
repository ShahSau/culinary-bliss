package middleware

import (
	"github.com/ShahSau/culinary-bliss/helpers"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware is a middleware that checks if the user is authenticated
func Authtication(c *gin.Context) {
	clientToken := c.Request.Header.Get("Authorization")
	if clientToken == "" {
		c.JSON(403, gin.H{"error": "No Authorization header provided"})
		c.Abort()
		return
	}

	claims, err := helpers.ValidateToken(clientToken)
	if err != nil {
		c.JSON(403, gin.H{"error": "Invalid token"})
		c.Abort()
		return
	}

	c.Set("email", claims.Email)
	c.Set("first_name", claims.First_name)
	c.Set("last_name", claims.Last_name)
	c.Set("user_id", claims.User_id)

	c.Next()

}
