package deployments_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestDeployments(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Deployments Suite")
}
