package matchers_test

import (
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"

	"testing"
)

func TestCollectors(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Test Matchers Suite")
}
