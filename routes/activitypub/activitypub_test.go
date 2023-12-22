package activitypub_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestActivityPub(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ActivityPub Router Suite")
}
