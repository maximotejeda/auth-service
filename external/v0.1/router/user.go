package router

import "github.com/gin-gonic/gin"

func UserAddRoutes(r *gin.Engine) {
	v01 := r.Group("/v0.1")
	user := v01.Group("/user")
	{
		user.GET("/info", nil)
		user.POST("/interaction", nil)
	}
}
