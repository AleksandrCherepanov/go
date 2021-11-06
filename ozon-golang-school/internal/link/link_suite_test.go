package link_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestLink(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Link Suite")
}
