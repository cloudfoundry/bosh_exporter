package deployments_test

import (
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"

	"testing"
)

func TestDeployments(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Deployments Suite")
}
