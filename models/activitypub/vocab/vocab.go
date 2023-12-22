// The vocab package is just a collection of consts of ActivityStreams (and extensions) terms, to get rid of magic strings.
//
// This package only includes vocab actually used by Korora, and is not suitable as a general purpose library.
package vocab

const (
	// core vocab
	Context = "@context"
	Type    = "type"
	Id      = "id"

	// activity types
	Activity = "Activity"
	Accept   = "Accept"

	// actor types
	Application  = "Application"
	Group        = "Group"
	Organization = "Organization"
	Person       = "Person"
	Service      = "Service"
	KororaHuddle = "korora:Huddle"

	// actor collections
	Inbox     = "inbox"
	Outbox    = "outbox"
	Following = "following"
	Followers = "followers"
	Liked     = "liked"

	// actor properties
	PreferredUsername = "preferredUsername"
	PublicKey         = "publicKey"
	Owner             = "owner"
	PublicKeyPem      = "publicKeyPem"

	MainKey = "#main-key"
)
