package router

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	db "github.com/maximotejeda/auth-service/external/v0.1/database"
	"github.com/maximotejeda/auth-service/external/v0.1/helper"
)

// AdminAddRoutes Add routes to predefined endpoints
func AdminAddRoutes(r *gin.Engine) {
	v01 := r.Group("/v0.1")
	admin := v01.Group("/admin")
	validated := admin.Group("/")
	validated.Use(helper.Validated())
	validated.Use(helper.IsAdmin())
	{
		validated.GET("/users", adminListUsers)
		validated.GET("/recover", adminInitiateUserRecovery)
		validated.POST("/activate", adminActivateAccount)
		validated.POST("/register", adminExternalAccountRegister)
		validated.POST("/ban", adminDisableAccount)
	}
}

func adminListUsers(c *gin.Context) {
	//TODO
	users, err := db.AdminGetUsers()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"users": users})
}

func adminInitiateUserRecovery(c *gin.Context) {
	//TODO
	data := c.Query("data")
	if data == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "account information needed"})
		return
	}
	db.NewRecoverAccount(data)
	c.JSON(http.StatusCreated, gin.H{"status": "if data exist an email will be sent in the next 10 minutes"})
}

func adminExternalAccountRegister(c *gin.Context) {
	//TODO
	input := db.Register{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// User active or inactive for default
	user := db.User{Active: true, UserLogin: input.UserLogin, UserInfo: input.UserInfo}
	err := db.CreateAccount(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"status": "created"})
}

func adminDisableAccount(c *gin.Context) {
	//Can be username or email
	data := c.Query("data")
	if data == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "account information needed"})
		return
	}
	err := db.AdminBanUser(data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"status": data + " Banned"})
}

func adminActivateAccount(c *gin.Context) {
	data := c.Query("data")
	if data == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "account information needed"})
		return
	}
	user := db.User{}
	if strings.Contains(data, "@") {
		user.Email = data
	} else {
		user.Username = data
	}
	db.ActivateAccount(&user)
	c.JSON(http.StatusCreated, gin.H{"status": data + " activated"})
}
