package activitypub_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/go-chi/chi/v5"
	"github.com/korora-social/korora/dao/daomocks"
	"github.com/korora-social/korora/models"
	"github.com/korora-social/korora/routes/activitypub"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("User-specific routes", func() {
	var (
		router  chi.Router
		userDao *daomocks.UserDao = nil
	)

	performRequest := func(req *http.Request) *http.Response {
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
		return rr.Result()
	}

	BeforeEach(func() {
		userDao = &daomocks.UserDao{
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

		Context("With an existing user's path", func() {
			It("Returns 200, with the user's activitypub representation", func() {
				req := httptest.NewRequest("GET", "/users/test1", nil)
				resp := performRequest(req)
				Expect(resp.StatusCode).To(Equal(200))

				unwrapped := map[string]interface{}{}
				Expect(json.NewDecoder(resp.Body).Decode(&unwrapped)).To(Succeed())
				Expect(unwrapped).To(HaveKey("@context"))
				Expect(unwrapped).To(HaveKeyWithValue("type", "Person"))
				Expect(unwrapped).To(HaveKeyWithValue("id", "https://foo.bar/ap/users/test1"))
			})
		})
	})
})
