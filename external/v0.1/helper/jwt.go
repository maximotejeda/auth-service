package helper

import (
	"crypto/rsa"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JWT struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

func NewJWT() JWT {
	priv, pub := MyGenerateKeys()
	return JWT{
		privateKey: priv,
		publicKey:  pub,
	}
}

func (j JWT) Create(ttl time.Duration, content interface{}) (token string, err error) {

	now := time.Now().UTC()

	claims := make(jwt.MapClaims)

	claims["dat"] = content             // Data we expect to have on JWT
	claims["exp"] = now.Add(ttl).Unix() // Expiration time after which the token ios invalid
	claims["iat"] = now.Unix()          // The time at wwhen the token was issued
	claims["nbf"] = now.Unix()          // The time before which the token must be disregarded.

	token, err = jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(j.privateKey)
	if err != nil {
		log.Printf("creating token %s", err.Error())
	}

	return token, nil
}

// Validate token in the algorithm used
// working as expected
func (j JWT) Validate(token string) (interface{}, error) {
	fmt.Print(len(strings.Split(token, ".")))
	tok, err := jwt.Parse(token, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected method: %s", jwtToken.Header["alg"])
		}
		return j.publicKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("Validate parsing: %w", err)
	}

	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok || !tok.Valid {
		return nil, fmt.Errorf("validate: invalid")
	}

	return claims["dat"], nil
}
