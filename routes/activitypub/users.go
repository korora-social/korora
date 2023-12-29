package activitypub

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/korora-social/korora/dao"
	serializer "github.com/korora-social/korora/models/activitypub"
	log "github.com/sirupsen/logrus"
)

func (r *Route) GetUser(rw http.ResponseWriter, req *http.Request) {
	username := chi.URLParam(req, "username")
	user, err := r.dao.User().GetByUsername(username)
	if err != nil {
		if err == dao.NotFound {
			rw.WriteHeader(http.StatusNotFound)
			return
		}

		log.WithError(err).WithField("username", username).Warn("Error trying to fetch user")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(rw).Encode(serializer.User(user))
	if err != nil {
		log.WithError(err).Warn("Couldn't encode user to activity pub")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
}
