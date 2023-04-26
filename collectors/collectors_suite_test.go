package collectors_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"testing"
)

const (
	testNamespace   = "test_exporter"
	testEnvironment = "test_environment"
	testBoshName    = "test_bosh_name"
	testBoshUUID    = "test_bosh_uuid"
)

func TestCollectors(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Collectors Suite")
}
