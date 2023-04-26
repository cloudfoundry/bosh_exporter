package filters_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/bosh-prometheus/bosh_exporter/filters"
)

var _ = Describe("CollectorsFilter", func() {
	var (
		err     error
		filters []string

		collectorsFilter *CollectorsFilter
	)

	JustBeforeEach(func() {
		collectorsFilter, err = NewCollectorsFilter(filters)
	})

	Describe("New", func() {
		Context("when filters are supported", func() {
			BeforeEach(func() {
				filters = []string{DeploymentsCollector, JobsCollector, ServiceDiscoveryCollector}
			})

			It("does not return an error", func() {
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("when filters are not supported", func() {
			BeforeEach(func() {
				filters = []string{DeploymentsCollector, JobsCollector, "Unknown"}
			})

			It("returns an error", func() {
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("collector filter `Unknown` is not supported"))
			})
		})

		Context("when a filter has leading and/or trailing whitespaces", func() {
			BeforeEach(func() {
				filters = []string{"   " + DeploymentsCollector + "  "}
			})

			It("returns an error", func() {
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})

	Describe("Enabled", func() {
		BeforeEach(func() {
			filters = []string{DeploymentsCollector}
		})

		Context("when collector is enabled", func() {
			It("returns true", func() {
				Expect(collectorsFilter.Enabled(DeploymentsCollector)).To(BeTrue())
			})
		})

		Context("when collector is not enabled", func() {
			It("returns false", func() {
				Expect(collectorsFilter.Enabled(JobsCollector)).To(BeFalse())
			})
		})

		Context("when there are no filters", func() {
			BeforeEach(func() {
				filters = []string{}
			})

			It("returns true", func() {
				Expect(collectorsFilter.Enabled(JobsCollector)).To(BeTrue())
			})
		})
	})
})
