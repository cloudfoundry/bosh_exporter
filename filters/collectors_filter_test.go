package filters_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/cloudfoundry-community/bosh_exporter/filters"
)

var _ = Describe("CollectorsFilter", func() {
	var (
		err    error
		filter []string

		collectorsFilter *CollectorsFilter
	)

	JustBeforeEach(func() {
		collectorsFilter, err = NewCollectorsFilter(filter)
	})

	Describe("New", func() {
		Context("when filters are supported", func() {
			BeforeEach(func() {
				filter = []string{DeploymentsCollector, JobsCollector}
			})

			It("does not return an error", func() {
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("when filters are not supported", func() {
			BeforeEach(func() {
				filter = []string{DeploymentsCollector, JobsCollector, "Unknown"}
			})

			It("returns an error", func() {
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("Collector filter `Unknown` is not supported"))
			})
		})
	})

	Describe("Enabled", func() {
		BeforeEach(func() {
			filter = []string{DeploymentsCollector}
		})

		Context("when event is enabled", func() {
			It("returns true", func() {
				Expect(collectorsFilter.Enabled(DeploymentsCollector)).To(BeTrue())
			})
		})

		Context("when event is not enabled", func() {
			It("returns false", func() {
				Expect(collectorsFilter.Enabled(JobsCollector)).To(BeFalse())
			})
		})

		Context("when there is no filter", func() {
			BeforeEach(func() {
				filter = []string{}
			})

			It("returns true", func() {
				Expect(collectorsFilter.Enabled(JobsCollector)).To(BeTrue())
			})
		})
	})
})
