package activitypub_test

import (
	"net/http"
	"net/http/httptest"

	"github.com/go-chi/chi/v5"
	"github.com/korora-social/korora/dao/user"
	"github.com/korora-social/korora/models"
	"github.com/korora-social/korora/routes/activitypub"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("User-specific routes", func() {
	var (
		router  chi.Router
		userDao *user.MockDao = nil
	)

	performRequest := func(req *http.Request) *http.Response {
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
		return rr.Result()
	}

	BeforeEach(func() {
		userDao = &user.MockDao{
			Users: []*models.User{
				&models.User{
					Username: "test1",
					Uri:      "https://foo.bar/ap/users/test1",
				},
				&models.User{
					Username: "test2",
					Uri:      "https://foo.bar/ap/users/test2",
				},
			},
		}
		apRoute := activitypub.New(userDao)
		router = chi.NewRouter()
		router.Route("/", apRoute.Routes)
	})

	Describe("User profile in activitypub", func() {
		Context("With a non-existent User path", func() {
			It("Returns a 404", func() {
				req := httptest.NewRequest("GET", "/users/foob", nil)
				resp := performRequest(req)
				Expect(resp.StatusCode).To(Equal(404))
			})
		})
	})
})
