package models_test

import (
	"crypto/x509"
	"encoding/pem"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/korora-social/korora/models"
)

var _ = Describe("User", func() {
	Describe("CreateUser", func() {
		It("Creates a basic User struct", func() {
			user := models.NewUser("foo", "https://foo.bar")
			Expect(user.Username).To(Equal("foo"))
			Expect(user.Uri).To(Equal("https://foo.bar/ap/users/foo"))
		})

		It("Creates a key pair for the user", func() {
			user := models.NewUser("foo", "https://foo.bar")
			Expect(user.PrivateKey).ToNot(BeEmpty())
			Expect(user.PublicKey).ToNot(BeEmpty())

			publicKeyPem := user.PublicKeyPem()

			// check that public key from the PEM matches the private key on the user
			p, _ := pem.Decode(publicKeyPem)
			pub, err := x509.ParsePKCS1PublicKey(p.Bytes)
			Expect(err).ToNot(HaveOccurred())
			Expect(pub).To(Equal(user.GetPrivateKey().Public()))
		})
	})
})
