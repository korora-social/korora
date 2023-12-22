package webfinger_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/korora-social/korora/routes/webfinger"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/korora-social/korora/dao/user"
	"github.com/korora-social/korora/models"
)

func TestWebfinger(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "WebFinger Suite")
}

type mockUserDao struct {
	user.Dao
	Users []*models.User
}

func (u *mockUserDao) GetByWebfinger(username, domain string) (*models.User, error) {
	for _, user := range u.Users {
		if user.Username == username {
			return user, nil
		}
	}

	return nil, fmt.Errorf("No user found")
}

var _ = Describe("Webfinger", func() {
	var router *webfinger.WebfingerRoute = nil
	var userDao *mockUserDao = nil

	BeforeEach(func() {
		userDao = &mockUserDao{
			Users: []*models.User{},
		}
		router = webfinger.New(userDao)
	})

	Describe("HTTP Get requests", func() {
		var respRecorder *httptest.ResponseRecorder = nil
		var request *http.Request = nil

		BeforeEach(func() {
			respRecorder = httptest.NewRecorder()
			request = httptest.NewRequest("GET", "/.well-known/webfinger?resource=acct:foo@bar.com", nil)
		})

		Context("Without any users in the database", func() {
			It("Should return 404", func() {
				router.ServeHTTP(respRecorder, request)

				result := respRecorder.Result()
				Expect(result.StatusCode).To(Equal(404))
			})
		})

		Context("With the given user in the database", func() {
			It("Should return a 200, and link to their ActivityPub url", func() {
				userDao.Users = []*models.User{
					&models.User{
						Username: "foo",
						Uri:      "https://bar.com/users/foo",
					},
				}

				router.ServeHTTP(respRecorder, request)

				result := respRecorder.Result()
				Expect(result.StatusCode).To(Equal(200))

				expected := webfinger.Response{
					Subject: "acct:foo@bar.com",
					Links: []webfinger.Link{{
						Rel:  "self",
						Type: "application/activity+json",
						Href: "https://bar.com/users/foo",
					}},
				}
				var actual webfinger.Response

				Expect(json.NewDecoder(result.Body).Decode(&actual)).Should(Succeed())
				Expect(actual).To(Equal(expected))
			})
		})
	})
})
