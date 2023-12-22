package activitypub_test

import (
	"github.com/korora-social/korora/models"
	"github.com/korora-social/korora/models/activitypub"
	"github.com/korora-social/korora/models/ld"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("User Serializer", func() {
	Context("With a generic user", func() {
		var user *models.User

		BeforeEach(func() {
			user = models.NewUser("foo", "https://foo.bar")
		})

		It("Builds a basic ActivityPub Person", func() {
			document := activitypub.User(user)
			Expect(document).To(HaveKeyWithValue("@context", ld.BaseContext))
			Expect(document).To(HaveKeyWithValue("type", "Person"))
			Expect(document).To(HaveKeyWithValue("id", user.Uri))
		})
	})
})
