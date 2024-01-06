package models

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"net/url"
	"path"
)

const (
	UserKeyBits  int    = 2048
	UserBasePath string = "ap/users"
)

var (
	ErrKeyPresent error = errors.New("Key pair already set")
)

type User struct {
	Username string `sq:"username"`
	Uri      string `sq:"uri"`

	PrivateKey []byte `sq:"private_key"`
	PublicKey  []byte `sq:"public_key"`
}

func NewUser(username, baseUrl string) *User {
	url, _ := url.Parse(baseUrl)
	url.Path = path.Join(url.Path, UserBasePath, username)

	user := &User{
		Username: username,
		Uri:      url.String(),
	}

	user.GenerateKeys()

	return user
}

func (u *User) GenerateKeys() error {
	if len(u.PrivateKey) > 0 {
		return ErrKeyPresent
	}

	// create user's public and private keys
	key, err := rsa.GenerateKey(rand.Reader, UserKeyBits)
	if err != nil {
		return err
	}

	u.PrivateKey = x509.MarshalPKCS1PrivateKey(key)
	u.PublicKey = x509.MarshalPKCS1PublicKey(key.Public().(*rsa.PublicKey))

	return nil
}

func (u *User) PublicKeyPem() []byte {
	return pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: u.PublicKey,
	})
}

func (u *User) GetPrivateKey() *rsa.PrivateKey {
	key, err := x509.ParsePKCS1PrivateKey(u.PrivateKey)
	if err != nil {
		panic(err)
	}

	return key
}

// append path parts to a user's Uri
func (user *User) Path(parts ...string) string {
	u, _ := url.Parse(user.Uri)
	u.Path = path.Join(append([]string{u.Path}, parts...)...)
	return u.String()
}
