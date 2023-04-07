package helper

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type ValidatedRequest struct {
	Authorization string
}

func MyGenerateKeys() (priv *rsa.PrivateKey, pub *rsa.PublicKey) {

	priv, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		fmt.Printf("%v", err)
	}
	//fmt.Print(priv)
	pub = &priv.PublicKey

	priv.Validate()
	return priv, pub
}

// Validate token header to manage auth
func Validated() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		s := strings.Replace(token, "Bearer ", "", 1)
		_, err := GlobalKeys.Validate(s)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"err": err.Error()})
			c.Abort()
			return
		} else {
			c.Next()
		}
	}
}
