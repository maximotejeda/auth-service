package router

import (
	"github.com/gin-gonic/gin"
	"github.com/maximotejeda/auth-service/external/v0.1/database"
	"github.com/maximotejeda/auth-service/external/v0.1/helper"
)

var R *gin.Engine
var J *helper.JWT = helper.GlobalKeys
var db = database.DB

// Return a new Engine
func NewRouter() *gin.Engine {
	if R == nil {
		R = gin.New()
		R.Use(gin.Logger())
	}
	return R
}
