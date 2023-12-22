package models

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"net/url"
	"path"
)

const (
	UserKeyBits  int    = 2048
	UserBasePath string = "ap/users"
)

type User struct {
	Username string
	Uri      string

	PrivateKey []byte
	PublicKey  []byte
}

func NewUser(username, baseUrl string) *User {
	url, _ := url.Parse(baseUrl)
	url.Path = path.Join(url.Path, UserBasePath, username)

	user := &User{
		Username: username,
		Uri:      url.String(),
	}

	// create user's public and private keys
	key, err := rsa.GenerateKey(rand.Reader, UserKeyBits)
	if err != nil {
		panic(err)
	}

	user.PrivateKey = x509.MarshalPKCS1PrivateKey(key)
	user.PublicKey = x509.MarshalPKCS1PublicKey(key.Public().(*rsa.PublicKey))

	return user
}

func (u *User) PublicKeyPem() string {
	return string(pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: u.PublicKey,
	}))
}

func (u *User) GetPrivateKey() *rsa.PrivateKey {
	key, err := x509.ParsePKCS1PrivateKey(u.PrivateKey)
	if err != nil {
		panic(err)
	}

	return key
}
