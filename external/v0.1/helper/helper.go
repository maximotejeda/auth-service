package helper

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
)

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
