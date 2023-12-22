package activitypub

import (
	"github.com/go-chi/chi/v5"
	"github.com/korora-social/korora/dao/user"
)

type Route struct {
	users user.Dao
}

func New(users user.Dao) *Route {
	return &Route{
		users: users,
	}
}

func (r *Route) Routes(rtr chi.Router) {
	rtr.Route("/users", func(rtr chi.Router) {
		rtr.Route("/{username}", func(rtr chi.Router) {
			rtr.Get("/", r.GetUser)
		})
	})
}
