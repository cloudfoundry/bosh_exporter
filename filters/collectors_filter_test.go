package filters_test

import (
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"

	"github.com/cloudfoundry/bosh_exporter/filters"
)

var _ = ginkgo.Describe("CollectorsFilter", func() {
	var (
		err          error
		filtersArray []string

		collectorsFilter *filters.CollectorsFilter
	)

	ginkgo.JustBeforeEach(func() {
		collectorsFilter, err = filters.NewCollectorsFilter(filtersArray)
	})

	ginkgo.Describe("New", func() {
		ginkgo.Context("when filters are supported", func() {
			ginkgo.BeforeEach(func() {
				filtersArray = []string{filters.DeploymentsCollector, filters.JobsCollector, filters.ServiceDiscoveryCollector}
			})

			ginkgo.It("does not return an error", func() {
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
			})
		})

		ginkgo.Context("when filters are not supported", func() {
			ginkgo.BeforeEach(func() {
				filtersArray = []string{filters.DeploymentsCollector, filters.JobsCollector, "Unknown"}
			})

			ginkgo.It("returns an error", func() {
				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(err.Error()).To(gomega.Equal("collector filter `Unknown` is not supported"))
			})
		})

		ginkgo.Context("when a filter has leading and/or trailing whitespaces", func() {
			ginkgo.BeforeEach(func() {
				filtersArray = []string{"   " + filters.DeploymentsCollector + "  "}
			})

			ginkgo.It("returns an error", func() {
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
			})
		})
	})

	ginkgo.Describe("Enabled", func() {
		ginkgo.BeforeEach(func() {
			filtersArray = []string{filters.DeploymentsCollector}
		})

		ginkgo.Context("when collector is enabled", func() {
			ginkgo.It("returns true", func() {
				gomega.Expect(collectorsFilter.Enabled(filters.DeploymentsCollector)).To(gomega.BeTrue())
			})
		})

		ginkgo.Context("when collector is not enabled", func() {
			ginkgo.It("returns false", func() {
				gomega.Expect(collectorsFilter.Enabled(filters.JobsCollector)).To(gomega.BeFalse())
			})
		})

		ginkgo.Context("when there are no filters", func() {
			ginkgo.BeforeEach(func() {
				filtersArray = []string{}
			})

			ginkgo.It("returns true", func() {
				gomega.Expect(collectorsFilter.Enabled(filters.JobsCollector)).To(gomega.BeTrue())
			})
		})
	})
})
