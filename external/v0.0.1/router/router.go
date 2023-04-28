package router

import (
	"errors"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/maximotejeda/helpers/jwts"
)

var R *gin.Engine
var J *jwts.JWT = &jwts.JWT{}

func init() {
	_, err := os.Lstat(jwts.PrivateKeyDir)

	if errors.Is(err, os.ErrNotExist) {
		// this means the keys are not writed to disk so ill create them
		fmt.Println("\nKeys does not exist writing them", err)
		J.New() // write to disk on first run
	} else if err == nil {
		J.ReadFromDisk() // read keys if other instance is running

		// TODO priodically check if keys changed
		// im planning on changing keys each day or half a day
	}

}

// Return a new Engine
func NewRouter() *gin.Engine {
	if R == nil {
		R = gin.New()
		R.Use(gin.Logger())
	}
	return R
}
