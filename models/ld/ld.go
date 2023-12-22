// The ld package is used for normalizing json-ld documents
package ld

import (
	"github.com/piprate/json-gold/ld"
)

type Processor struct {
	proc *ld.JsonLdProcessor
	opts *ld.JsonLdOptions
}

var (
	BaseContext []interface{} = []interface{}{
		"https://www.w3.org/ns/activitystreams",
		"https://w3id.org/security/v1",
		map[string]interface{}{
			"toot":                      "http://joinmastodon.org/ns#",
			"manuallyApprovesFollowers": "as:manuallyApprovesFollowers",
			"alsoKnownAs":               "as:alsoKnownAs",
			"attachment":                "as:attachment",
		},
	}

	instance *Processor = nil
)

func Get() *Processor {
	if instance == nil {
		instance = &Processor{
			proc: ld.NewJsonLdProcessor(),
			opts: ld.NewJsonLdOptions(""),
		}
	}

	return instance
}

func (p *Processor) Normalize(document map[string]interface{}) (map[string]interface{}, error) {
	if p == nil {
		p = Get()
	}

	expanded, err := p.proc.Expand(document, p.opts)
	if err != nil {
		return nil, err
	}

	return p.proc.Compact(expanded, BaseContext, p.opts)
}
