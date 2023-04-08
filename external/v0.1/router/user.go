package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/maximotejeda/auth-service/external/v0.1/helper"
)

func UserAddRoutes(r *gin.Engine) {
	v01 := r.Group("/v0.1")
	user := v01.Group("/user")
	user.Use(helper.Validated())
	{
		user.GET("/info", userInfo)
		user.GET("/interaction", userInteractions)
	}
}

// userInfo returns information about the user
func userInfo(c *gin.Context) {
	username, _ := c.Get("username")
	email, _ := c.Get("email")
	rol, _ := c.Get("rol")
	c.JSON(http.StatusOK, gin.H{"username": username, "email": email, "rol": rol})
}

func userInteractions(c *gin.Context) {
	loged, _ := c.Get("loged")
	c.JSON(http.StatusOK, gin.H{"last Loged": loged})
}
