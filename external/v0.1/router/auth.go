package router

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/maximotejeda/auth-service/external/v0.1/helper"
)

func AuthAddRoutes(r *gin.Engine) {
	v01 := r.Group("/v0.1")
	auth := v01.Group("/auth")
	{
		auth.POST("/login", login)
		auth.POST("/register", nil)
		auth.POST("/validate", nil)
		auth.POST("/refresh", nil)
		auth.POST("/signout", nil)
	}
}

func login(c *gin.Context) {
	input := userLogin{}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	j := helper.NewJWT()

	s, err := j.Create(time.Second*50, input)
	if err != nil {
		log.Print("Create token: ", err.Error())
	}

	_, err = j.Validate(s)
	if err != nil {
		fmt.Printf("validating err: %v", err)
	}

	c.JSON(http.StatusOK, gin.H{"message": "validated"})
}
