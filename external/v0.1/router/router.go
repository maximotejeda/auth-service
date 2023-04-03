package router

import (
	"github.com/gin-gonic/gin"
)

var R *gin.Engine

// Return a new Engine
func NewRouter() *gin.Engine {
	if R == nil {
		R = gin.Default()
	}
	return R
}
