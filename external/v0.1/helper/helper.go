package helper

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type ValidatedRequest struct {
	Authorization string
}

// Validate token header to manage auth
func Validated() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		s := strings.Replace(token, "Bearer ", "", 1)
		_, err := GlobalKeys.Validate(s)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"err": err.Error()})
			c.Abort()
			return
		} else {
			c.Next()
		}
	}
}
