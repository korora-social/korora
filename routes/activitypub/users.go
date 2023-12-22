package activitypub

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/korora-social/korora/dao"
	log "github.com/sirupsen/logrus"
)

func (r *Route) GetUser(rw http.ResponseWriter, req *http.Request) {
	username := chi.URLParam(req, "username")
	user, err := r.users.GetByUsername(username)
	if err != nil {
		if err == dao.NotFound {
			rw.WriteHeader(http.StatusNotFound)
			return
		}

		log.WithError(err).WithField("username", username).Warn("Error trying to fetch user")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(rw).Encode(user.AsActivityPub())
	if err != nil {
		log.WithError(err).Warn("Couldn't encode user to activity pub")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
}
