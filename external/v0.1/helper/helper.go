package helper

import (
	"net/http"
	"strings"

	//"github.com/dgrijalva/jwt-go"
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
		claims, err := GlobalKeys.Validate(s)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"err": err.Error()})
			c.Abort()
			return
		} else {

			// add information to request to reduce information decrypt from JWT
			username, _ := claims["username"].(string)
			email, _ := claims["email"].(string)
			rol, _ := claims["rol"].(string)
			loged, _ := claims["loged"].(string)
			c.Set("username", username)
			c.Set("email", email)
			c.Set("rol", rol)
			c.Set("loged", loged)
		}
		c.Next()

	}
}

// Function to verify rol of a user and  give auth to certasins functions
func IsAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenstr := c.GetHeader("Authorization")
		s := strings.Replace(tokenstr, "Bearer ", "", 1)
		params, err := GlobalKeys.Validate(s)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"err": err.Error()})
			c.Abort()
			return
		}
		rol := params["rol"]
		rol, _ = rol.(string)
		if rol != "admin" {
			c.JSON(http.StatusUnauthorized, gin.H{"err": "Only admin allowed here"})
			c.Abort()
			return
		}
		c.Next()
	}
}
