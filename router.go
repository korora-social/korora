package korora

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/korora-social/korora/routes/webfinger"
)

func (k *Korora) Router() http.Handler {
	rtr := chi.NewRouter()

	rtr.Route("/.well-known", func(r chi.Router) {
		wf := webfinger.New(k.usersDao)
		wf.AddRoutes(rtr)
	})

	return rtr
}
