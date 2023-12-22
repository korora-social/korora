package webfinger

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/korora-social/korora/dao/user"
	log "github.com/sirupsen/logrus"
)

type Route struct {
	users user.Dao
}

type Response struct {
	Subject string `json:"subject"`
	Links   []Link `json:"links"`
}

type Link struct {
	Rel  string `json:"rel"`
	Type string `json:"type"`
	Href string `json:"href"`
}

func New(users user.Dao) *Route {
	return &Route{
		users: users,
	}
}

func (wf *Route) AddRoutes(rtr chi.Router) {
	rtr.Get("/webfinger", wf.Get)
}

func (wf *Route) Get(rw http.ResponseWriter, req *http.Request) {
	query := req.URL.Query()

	if !query.Has("resource") {
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	resource := query.Get("resource")
	userQuery, isAcct := strings.CutPrefix(resource, "acct:")
	if !isAcct {
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	username, _, hasAt := strings.Cut(userQuery, "@")
	if !hasAt {
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	user, err := wf.users.GetByUsername(username)
	if err != nil {
		log.WithError(err).WithFields(log.Fields{
			"username": username,
		}).Warn("Webfinger lookup for non-existent user")

		rw.WriteHeader(http.StatusNotFound)
		return
	}

	response := Response{
		Subject: resource,
		Links: []Link{{
			Rel:  "self",
			Type: "application/activity+json",
			Href: user.Uri,
		}},
	}

	err = json.NewEncoder(rw).Encode(&response)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
}
