package router

import "github.com/gin-gonic/gin"

func AdminAddRoutes(r *gin.Engine) {
	admin := r.Group("/admin")
	{
		admin.GET("/user", nil)
		admin.GET("/recover", nil)
		admin.POST("/validate", nil)
		admin.GET("/register", nil)
		admin.GET("/ban", nil)
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
