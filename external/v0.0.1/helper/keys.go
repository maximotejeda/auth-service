package helper

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
)

// function to generate key pairs from random 4096 bits
func MyGenerateKeys() (priv *rsa.PrivateKey, pub *rsa.PublicKey) {

	priv, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		fmt.Printf("%v", err)
	}
	//fmt.Print(priv)
	pub = &priv.PublicKey

	priv.Validate()
	return priv, pub
}

// Parse private keys from pem
func parsePrivateKeyFromPemStr(privPEM string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(privPEM))
	if block == nil {
		return nil, errors.New("failed to parse PEM block contianing key")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return priv, nil
}

// Encode the key to be able to save to disk and reuse or expor to configmap
func exportRSAPrivateKeyAsPemStr(privKey *rsa.PrivateKey) string {
	privKeyBytes := x509.MarshalPKCS1PrivateKey(privKey)
	privKey_pem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: privKeyBytes,
		},
	)
	return string(privKey_pem)
}

// Parse private keys from pem
func ParsePublicKeyFromPemStr(pubPEM string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pubPEM))
	if block == nil {
		return nil, errors.New("failed to parse PEM block from public key ")
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	switch pub := pub.(type) {
	case *rsa.PublicKey:
		return pub, nil
	default:
		break
	}
	return nil, errors.New("key type is not rsa")
}

// Encode the key to be able to save to disk and reuse or expor to configmap
func exportRSAPublicKeyAsPemStr(pubKey *rsa.PublicKey) (string, error) {
	pubKeyBytes, err := x509.MarshalPKIXPublicKey(pubKey)
	if err != nil {
		return "", err
	}
	pubKeyPem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: pubKeyBytes,
		},
	)
	return string(pubKeyPem), nil
}
