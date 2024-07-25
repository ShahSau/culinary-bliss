package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware is a middleware that checks if the user is authenticated
func Authtication(c *gin.Context) {
	fmt.Println("AuthMiddleware")
}
