package activitypub

import (
	"fmt"

	"github.com/korora-social/korora/models"
	"github.com/korora-social/korora/models/activitypub/vocab"
	"github.com/korora-social/korora/models/ld"
)

func User(user *models.User) map[string]interface{} {
	document := map[string]interface{}{
		vocab.Context: ld.BaseContext,
		vocab.Type:    vocab.Person,
		vocab.Id:      user.Uri,

		vocab.PreferredUsername: user.Username,

		vocab.Inbox:  user.Path("inbox"),
		vocab.Outbox: user.Path("outbox"),

		vocab.PublicKey: PublicKey(user.PublicKeyPem(), user.Uri, vocab.MainKey),
	}

	return document
}

func PublicKey(pem []byte, uri, id string) map[string]interface{} {
	return map[string]interface{}{
		vocab.Id:           fmt.Sprintf("%s%s", uri, id),
		vocab.Owner:        uri,
		vocab.PublicKeyPem: string(pem),
	}
}
