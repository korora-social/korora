package webfinger

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/korora-social/korora/dao/user"
	log "github.com/sirupsen/logrus"
)

type WebfingerRoute struct {
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

func New(users user.Dao) *WebfingerRoute {
	return &WebfingerRoute{
		users: users,
	}
}

func (wf *WebfingerRoute) Router() http.Handler {
	return wf
}

func (wf *WebfingerRoute) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
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

	username, domain, hasAt := strings.Cut(userQuery, "@")
	if !hasAt {
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	user, err := wf.users.GetByWebfinger(username, domain)
	if err != nil {
		log.WithError(err).WithFields(log.Fields{
			"username": username,
			"domain":   domain,
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
