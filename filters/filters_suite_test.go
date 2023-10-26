package filters_test

import (
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"

	"testing"
)

func TestFilters(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Filters Suite")
}
