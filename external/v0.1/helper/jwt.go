package helper

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	//"github.com/dgrijalva/jwt-go"
	"github.com/golang-jwt/jwt/v5"
)

type JWT struct {
	privateKey    *rsa.PrivateKey
	publicKey     *rsa.PublicKey
	privatePemStr string
	PublicPemStr  string
}

var (
	GlobalKeys    *JWT = &JWT{}
	actualDir, _       = os.Getwd()
	keys               = os.Getenv("KEYSDIR")
	priv               = os.Getenv("PRIVATEKEYNAME")
	pub                = os.Getenv("PUBLICKEYNAME")
	privateKeyDir      = path.Join(actualDir, keys, priv)
	publicKeyDir       = path.Join(actualDir, keys, pub)
	tokenTTL           = os.Getenv("TOKENTTLS")
)

// Init Check if theres other keys created
// if others keys exists dont recreate them just copy them and use the context
func init() {
	_, err := os.Lstat(privateKeyDir)

	if errors.Is(err, os.ErrNotExist) {
		// this means the keys are not writed to disk so ill create them
		fmt.Println("\nKeys does not exist writing them", err)
		GlobalKeys.New() // write to disk on first run
	} else if err == nil {
		GlobalKeys.readFromDisk() // read keys if other instance is running

		// TODO priodically check if keys changed
		// im planning on changing keys each day or half a day
	}
}

func NewJWT() *JWT {
	priv, pub := MyGenerateKeys()
	return &JWT{
		privateKey: priv,
		publicKey:  pub,
	}
}
func (j *JWT) New() {
	if j.privateKey == nil {
		j.Renew()
	}
}

// Renew rotate the keys creating new keys in case of existence
func (j *JWT) Renew() {
	priv, pub := MyGenerateKeys()
	j.privateKey = priv
	j.publicKey = pub
	// need to try to convert to string
	j.privatePemStr = exportRSAPrivateKeyAsPemStr(priv)
	j.PublicPemStr, _ = exportRSAPublicKeyAsPemStr(pub)
	j.writeToDisk()
	//fmt.Printf("\nPrivatekeyPem: %s\n\n PublicKeyPem: %s", j.privatePemStr, j.publicPemStr)
}

// writeToDisk every time new keys are issued write them to disk overwriting the actuals
func (j *JWT) writeToDisk() {
	err := os.Mkdir("keys", 0770)
	if err != nil {
		fmt.Printf("\n%v\n", err)
	}

	err = os.WriteFile(publicKeyDir, []byte(GlobalKeys.PublicPemStr), 0644)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	err = os.WriteFile(privateKeyDir, []byte(GlobalKeys.privatePemStr), 0644)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
}

func (j *JWT) readFromDisk() {
	privateBytes, err := os.ReadFile(privateKeyDir)
	if err != nil {
		panic(err)
	}
	publicBytes, err := os.ReadFile(publicKeyDir)
	if err != nil {
		panic(err)
	}
	// assig string to J
	// mostly i intend to use keys as a way of replicate the service throug multiples pods
	// share the same keys with multiples pods if needed to replicate
	GlobalKeys.privatePemStr = string(privateBytes)
	GlobalKeys.PublicPemStr = string(publicBytes)
	GlobalKeys.privateKey, err = parsePrivateKeyFromPemStr(GlobalKeys.privatePemStr)
	if err != nil {
		panic(err)
	}
	GlobalKeys.publicKey, err = ParsePublicKeyFromPemStr(GlobalKeys.PublicPemStr)
	if err != nil {
		panic(err)
	}
}

func (j *JWT) Create(content interface{}) (token string, err error) {
	if j == nil || j.privateKey == nil {
		log.Print("nil struct pointer")
		return "", fmt.Errorf("nil pointer struct %v", j)
	}
	tokenTimeToLiveInt, err := strconv.Atoi(tokenTTL)
	if err != nil {
		tokenTimeToLiveInt = 900
	}
	tokenTimeToLive := time.Second * time.Duration(tokenTimeToLiveInt)
	now := time.Now().UTC()

	claims := make(jwt.MapClaims)

	claims["dat"] = content                         // Data we expect to have on JWT
	claims["exp"] = now.Add(tokenTimeToLive).Unix() // Expiration time after which the token ios invalid
	claims["iat"] = now.Unix()                      // The time at wwhen the token was issued
	claims["nbf"] = now.Unix()                      // The time before which the token must be disregarded.

	token, err = jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(j.privateKey)
	if err != nil {
		log.Printf("creating token %s", err.Error())
	}

	return token, nil
}

// Validate token in the algorithm used
// working as expected
// return claims["dat"]
func (j *JWT) Validate(token string) (map[string]interface{}, error) {
	if j == nil || j.privateKey == nil {
		log.Print("nil struct pointer")
		return nil, fmt.Errorf("nil pointer struct %v", j)
	}

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
	dat, ok := claims["dat"].(map[string]interface{})
	if !ok {
		return map[string]interface{}{}, nil
	}
	return dat, nil
}

// RefreshToken: check header token validates it and check if is expired
// if expiration date is less than x issue a new
// Todo Set time for token and expiration timeout
func (j *JWT) RefreshToken(tokenStr string) (string, error) {
	_, err2 := j.Validate(tokenStr)
	var expNum int64
	switch {

	case err2 != nil && errors.Is(err2, jwt.ErrTokenExpired):
		token, err1 := jwt.Parse(tokenStr, nil)
		if token == nil {
			return "", err1
		}
		claims, _ := token.Claims.(jwt.MapClaims)
		// When the token expired?
		exp := claims["exp"]
		switch t := exp.(type) {
		case int:
			expNum = int64(t)
		case int64:
			expNum = t
		case float64:
			expNum = int64(t)
		default:
			fmt.Printf("The value is %v", t)
		}

		now := time.Now().Unix()
		tokenTTLInt, _ := strconv.Atoi(tokenTTL)
		if tokenTTLInt <= 0 {
			tokenTTLInt = 60
		}
		if (now-expNum) <= int64(tokenTTLInt) && (now-expNum) > 0 {
			//TODO modify the durarions
			newToken, err := j.Create(claims["dat"])
			if err != nil {
				return "", err
			} else {
				return newToken, err
			}
		} else {
			return "", fmt.Errorf("need new token session expired long")
		}
	case err2 != nil && !strings.Contains(err2.Error(), "Token is expired"):
		return "", err2
	}
	fmt.Println("Token still valid needs to wait for expiration")
	return tokenStr, nil
}
