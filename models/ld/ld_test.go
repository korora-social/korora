package ld_test

import (
	"github.com/korora-social/korora/models/ld"

	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestLd(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "JSON-LD Test Suite")
}

var _ = Describe("JSON-LD normalization", func() {
	Context("With a Mastodon format profile document", func() {
		processor := ld.Get()

		document := map[string]interface{}{
			"@context": []interface{}{
				"https://www.w3.org/ns/activitystreams",
				"https://w3id.org/security/v1",
				map[string]interface{}{
					"manuallyApprovesFollowers": "as:manuallyApprovesFollowers",
					"toot":                      "http://joinmastodon.org/ns#",
					"featured": map[string]interface{}{
						"@id":   "toot:featured",
						"@type": "@id",
					},
					"featuredTags": map[string]interface{}{
						"@id":   "toot:featuredTags",
						"@type": "@id",
					},
					"alsoKnownAs": map[string]interface{}{
						"@id":   "as:alsoKnownAs",
						"@type": "@id",
					},
					"movedTo": map[string]interface{}{
						"@id":   "as:movedTo",
						"@type": "@id",
					},
					"schema":           "http://schema.org#",
					"PropertyValue":    "schema:PropertyValue",
					"value":            "schema:value",
					"discoverable":     "toot:discoverable",
					"Device":           "toot:Device",
					"Ed25519Signature": "toot:Ed25519Signature",
					"Ed25519Key":       "toot:Ed25519Key",
					"Curve25519Key":    "toot:Curve25519Key",
					"EncryptedMessage": "toot:EncryptedMessage",
					"publicKeyBase64":  "toot:publicKeyBase64",
					"deviceId":         "toot:deviceId",
					"claim": map[string]interface{}{
						"@type": "@id",
						"@id":   "toot:claim",
					},
					"fingerprintKey": map[string]interface{}{
						"@type": "@id",
						"@id":   "toot:fingerprintKey",
					},
					"identityKey": map[string]interface{}{
						"@type": "@id",
						"@id":   "@toot:identityKey",
					},
					"devices": map[string]interface{}{
						"@type": "@id",
						"@id":   "toot:devices",
					},
					"messageFranking": "toot:messageFranking",
					"messageType":     "toot:messageType",
					"cipherText":      "toot:cipherText",
					"suspended":       "toot:suspended",
					"memorial":        "toot:memorial",
					"indexable":       "toot:indexable",
					"focalPoint": map[string]interface{}{
						"@container": "@list",
						"@id":        "toot:focalPoint",
					},
				},
			},
			"id":                        "https://foo.bar/users/fakeuser",
			"type":                      "Person",
			"following":                 "https://foo.bar/users/fakeuser/following",
			"followers":                 "https://foo.bar/users/fakeuser/followers",
			"inbox":                     "https://foo.bar/users/fakeuser/inbox",
			"outbox":                    "https://foo.bar/users/fakeuser/outbox",
			"featured":                  "https://foo.bar/users/fakeuser/collections/featured",
			"featuredTags":              "https://foo.bar/users/fakeuser/collections/tags",
			"preferredUsername":         "fakeuser",
			"name":                      "A Fake User",
			"summary":                   "<p>A biography for the fake user</p>",
			"url":                       "https://foo.bar/@fakeuser",
			"manuallyApprovesFollowers": false,
			"discoverable":              false,
			"indexable":                 false,
			"published":                 "2023-12-22T00:00:00Z",
			"memorial":                  false,
			"devices":                   "https://foo.bar/users/fakeuser/collections/devices",
			"alsoKnownAs": []string{
				"https://old.website/users/fakeuser",
			},
			"publicKey": map[string]interface{}{
				"id":           "https://foo.bar/users/fakeuser#main-key",
				"owner":        "https://foo.bar/users/fakeuser",
				"publicKeyPem": "-----BEGIN PUBLIC KEY-----\nfoobarbaz\n-----END PUBLIC KEY-----\n",
			},
			"tags": []interface{}{},
			"attachment": []map[string]interface{}{
				map[string]interface{}{
					"type":  "ProperyValue",
					"name":  "first property",
					"value": "first value",
				},
				map[string]interface{}{
					"type":  "PropertyValue",
					"name":  "another property",
					"value": "another value",
				},
			},
			"endpoints": map[string]interface{}{
				"sharedInbox": "https://foo.bar/inbox",
			},
			"icon": map[string]interface{}{
				"type":      "Image",
				"mediaType": "images/jpeg",
				"url":       "https://foo.bar/media/avatars/fakeuser.jpeg",
			},
			"image": map[string]interface{}{
				"type":      "Image",
				"mediaType": "image/png",
				"url":       "https://foo.bar/media/avatars/fakeuser-header.png",
			},
		}

		It("Should properly normalize the document", func() {
			document, err := processor.Normalize(document)
			Expect(err).ToNot(HaveOccurred())
			Expect(document).ToNot(BeEmpty())

			expectedKeys := []string{
				"@context", "id", "type", "following", "followers", "inbox", "outbox", "toot:featured", "toot:featuredTags", "preferredUsername", "name", "summary", "url", "manuallyApprovesFollowers", "toot:discoverable", "toot:indexable", "published", "toot:memorial", "toot:devices", "alsoKnownAs", "publicKey", "tags", "attachment", "icon", "image", "endpoints",
			}

			actualKeys := make([]string, len(document))
			i := 0
			for k := range document {
				actualKeys[i] = k
				i++
			}

			Expect(actualKeys).To(ConsistOf(expectedKeys))
		})
	})
})
