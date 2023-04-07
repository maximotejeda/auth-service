package router

import (
	"github.com/gin-gonic/gin"
	"github.com/maximotejeda/auth-service/external/v0.1/helper"
)

func UserAddRoutes(r *gin.Engine) {
	v01 := r.Group("/v0.1")
	user := v01.Group("/user")
	user.Use(helper.Validated())
	{
		user.GET("/info", nil)
		user.POST("/interaction", nil)
	}
}
