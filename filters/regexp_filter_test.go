package filters_test

import (
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"

	"github.com/bosh-prometheus/bosh_exporter/filters"
)

var _ = ginkgo.Describe("RegexpFilter", func() {
	var (
		err          error
		filtersArray []string

		regexpFilter *filters.RegexpFilter
	)

	ginkgo.JustBeforeEach(func() {
		regexpFilter, err = filters.NewRegexpFilter(filtersArray)
	})

	ginkgo.Describe("New", func() {
		ginkgo.Context("when filters compile", func() {
			ginkgo.BeforeEach(func() {
				filtersArray = []string{"bosh_exporter", "[a-z]+_collector"}
			})

			ginkgo.It("does not return an error", func() {
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
			})
		})

		ginkgo.Context("when filters does not compile", func() {
			ginkgo.BeforeEach(func() {
				filtersArray = []string{"[a-(z]+_exporter"}
			})

			ginkgo.It("returns an error", func() {
				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(err.Error()).To(gomega.Equal("error parsing regexp: invalid character class range: `a-(`"))
			})
		})
	})

	ginkgo.Describe("Enabled", func() {
		ginkgo.BeforeEach(func() {
			filtersArray = []string{"bosh_exporter", "[a-z]+_collector"}
		})

		ginkgo.Context("when there is a match", func() {
			ginkgo.It("returns true", func() {
				gomega.Expect(regexpFilter.Enabled("deployments_collector")).To(gomega.BeTrue())
			})
		})

		ginkgo.Context("when there is not a match", func() {
			ginkgo.It("returns false", func() {
				gomega.Expect(regexpFilter.Enabled("deployments_exporter")).To(gomega.BeFalse())
			})
		})

		ginkgo.Context("when there are no filters", func() {
			ginkgo.BeforeEach(func() {
				filtersArray = []string{}
			})

			ginkgo.It("returns true", func() {
				gomega.Expect(regexpFilter.Enabled("deployments_exporter")).To(gomega.BeTrue())
			})
		})
	})
})
