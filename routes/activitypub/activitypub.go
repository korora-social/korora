package activitypub

import (
	"github.com/go-chi/chi/v5"
	"github.com/korora-social/korora/dao"
)

type Route struct {
	dao dao.Dao
}

func New(dao dao.Dao) *Route {
	return &Route{
		dao: dao,
	}
}

func (r *Route) Routes(rtr chi.Router) {
	rtr.Route("/users", func(rtr chi.Router) {
		rtr.Route("/{username}", func(rtr chi.Router) {
			rtr.Get("/", r.GetUser)
		})
	})
}
