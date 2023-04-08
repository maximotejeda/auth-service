package router

import (
	"github.com/gin-gonic/gin"
	"github.com/maximotejeda/auth-service/external/v0.1/helper"
)

func AdminAddRoutes(r *gin.Engine) {
	v01 := r.Group("/v0.1")
	admin := v01.Group("/admin")
	validated := admin.Group("/")
	validated.Use(helper.Validated())
	validated.Use(helper.IsAdmin())
	{
		validated.GET("/users", nil)
		validated.GET("/recover", nil)
		validated.POST("/validate", nil)
		validated.GET("/register", nil)
		validated.GET("/ban", nil)
	}
}

func listUsers(c *gin.Context) {
	//TODO
}

func initiateUserRecovery(c *gin.Context) {
	//TODO
}

func externalAccountRegister(c *gin.Context) {
	//TODO
}

func disableAccount(c *gin.Context) {
	//TODO
}
